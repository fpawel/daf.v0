package main

import (
	"context"
	"fmt"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

func getComports() []string {
	ports, _ := comport.AvailablePorts()
	return ports
}

func comportIndex(portName string) int {
	ports, _ := comport.AvailablePorts()
	for i, s := range ports {
		if s == portName {
			return i
		}
	}
	return -1
}

func runMainWindow() error {

	app := walk.App()
	app.SetOrganizationName("analitpribor")
	app.SetProductName("EN8800-6408")
	settings := walk.NewIniFileSettings("settings.ini")
	if err := settings.Load(); err != nil {
		log.Fatal(err)
	}
	app.SetSettings(settings)

	getIniValue := func(key string) string {
		s, _ := settings.Get(key)
		return s
	}

	newComboBoxComport := func(comboBox **walk.ComboBox, key string) ComboBox {
		return ComboBox{
			AssignTo:     comboBox,
			Model:        getComports(),
			CurrentIndex: comportIndex(getIniValue(key)),
			OnMouseDown: func(_, _ int, _ walk.MouseButton) {
				cb := *comboBox
				n := cb.CurrentIndex()
				_ = cb.SetModel(getComports())
				_ = cb.SetCurrentIndex(n)
			},
			OnCurrentIndexChanged: func() {
				_ = settings.Put(key, (*comboBox).Text())
			},
		}
	}

	var (
		cbType,
		cbComportDaf,
		cbComportHart *walk.ComboBox
		nePgs      [4]*walk.NumberEdit
		tvProducts *walk.TableView

		pnWork       *walk.Composite
		neCmd, neArg *walk.NumberEdit
		pbCancelWork *walk.PushButton
	)

	var currentParty = data.LastParty()

	showErr := func(title, text string) {
		walk.MsgBox(mainWindow, title,
			text, walk.MsgBoxIconError|walk.MsgBoxOK)
	}

	saveParty := func() {
		if err := currentParty.Save(); err != nil {
			showErr("Ошибка данных", fmt.Sprintf("%v: %v", currentParty, err))
		}
	}

	var workStarted bool
	doWork := func(what string, work func() error) {
		if workStarted {
			panic("already started")
		}
		workStarted = true
		comportContext, cancelComport = context.WithCancel(context.Background())
		pnWork.SetVisible(false)
		portDaf.PortName = cbComportDaf.Text()
		portHart.PortName = cbComportHart.Text()
		pbCancelWork.SetVisible(true)
		go func() {
			err := work()
			_ = portHart.Port.Close()
			_ = portDaf.Port.Close()
			mainWindow.Synchronize(func() {
				workStarted = false
				pnWork.SetVisible(true)
				pbCancelWork.SetVisible(false)
				lastPartyProductsModel.setInterrogatePlace(-1)
				if err != nil {
					showErr(what, err.Error())
				}

			})
		}()
	}

	executeProductDialog := func() {
		n := tvProducts.CurrentIndex()
		if n < 0 || n >= len(lastPartyProductsModel.items) {
			return
		}
		p := *lastPartyProductsModel.items[n].Product
		cmd, err := runProductDialog(mainWindow, &p)
		if err != nil {
			showErr("Ошибка данных", err.Error())
			return
		}
		if cmd != walk.DlgCmdOK {
			return
		}
		if err := p.Save(); err != nil {
			showErr("Ошибка данных", err.Error())
			return
		}
		lastPartyProductsModel.items[n].Product = &p
		lastPartyProductsModel.PublishRowChanged(n)
	}

	newNumberEditPgs := func(n data.Gas) NumberEdit {
		return NumberEdit{
			Value:    currentParty.Pgs(n),
			AssignTo: &nePgs[n-1],
			MinValue: 0,
			Decimals: 2,
			OnValueChanged: func() {
				currentParty.SetPgs(n, nePgs[n-1].Value())
				saveParty()
			},
		}
	}

	if err := (MainWindow{
		AssignTo: &mainWindow,
		Title: fmt.Sprintf("ЭН8800-6408 Партия ДАФ-М №%d %s", currentParty.PartyID,
			currentParty.CreatedAt.Format("02.01.2006")),
		Name:       "MainWindow",
		Font:       Font{PointSize: 12, Family: "Segoe UI"},
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Size:       Size{800, 600},
		Layout:     VBox{},

		Children: []Widget{
			ScrollView{
				VerticalFixed: true,
				Layout:        HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &pbCancelWork,
						Text:     "Прервать",
						OnClicked: func() {
							cancelComport()
						},
					},

					Composite{
						AssignTo: &pnWork,
						Layout: HBox{
							MarginsZero: true,
						},
						Children: []Widget{
							SplitButton{
								Text: "Управление",
								MenuItems: []MenuItem{
									Action{
										Text: "Опрос",
										OnTriggered: func() {
											doWork("опрос", interrogateProducts)
										},
									},
									Action{
										Text: "Настройка ДАФ-М",
									},
								},
							},
							Label{Text: "Код:"},
							NumberEdit{AssignTo: &neCmd},
							Label{Text: "Аргумент:"},
							NumberEdit{AssignTo: &neArg, Decimals: 2, MinSize: Size{80, 0}},
							PushButton{Text: "Выполнить", OnClicked: func() {
								cmd := modbus.DevCmd(neCmd.Value())
								arg := neArg.Value()
								doWork(fmt.Sprintf("Оправка команды %d,%v", cmd, arg), func() error {
									return sendCmd(cmd, arg)
								})
							}},
						},
					},
				},
			},
			ScrollView{
				Layout: HBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					TableView{
						AssignTo:                 &tvProducts,
						NotSortableByHeaderClick: true,
						LastColumnStretched:      true,
						CheckBoxes:               true,
						Model:                    lastPartyProductsModel,

						Columns: []TableViewColumn{
							{Title: "Адрес", Width: 80},
							{Title: "Номер", Width: 100},
							{Title: "Концентрация", Width: 150},
							{Title: "Ток", Width: 100},
							{Title: "Порог 1", Width: 120},
							{Title: "Порог 2", Width: 120},
							{Title: "Связь"},
						},
						ContextMenuItems: []MenuItem{
							Action{
								Text: "Создать новую партию",
							},
							Action{
								Text: "Добавить прибор в партию",
								Shortcut: Shortcut{
									Key: walk.KeyInsert,
								},
								OnTriggered: func() {
									lastPartyProductsModel.addNewProduct()
								},
							},
							Action{
								Text: "Удалить прибор из партии",
							},
							Action{
								Text:        "Изменить адрес, номер прибора",
								OnTriggered: executeProductDialog,
							},
						},
					},
					ScrollView{
						HorizontalFixed: true,
						Layout:          VBox{},
						Children: []Widget{
							GroupBox{
								Title:  "COM порты",
								Layout: VBox{},
								Children: []Widget{
									Label{Text: "Стенд и приборы:"},
									newComboBoxComport(&cbComportDaf, "comport_products"),
									Label{Text: "HART модем:"},
									newComboBoxComport(&cbComportHart, "comport_hart"),
								},
							},
							Label{Text: "Исполнение:"},
							ComboBox{
								Model:         productTypes,
								AssignTo:      &cbType,
								DisplayMember: "Name",
								CurrentIndex:  indexOfProductTypeCode(currentParty.Type),
								OnCurrentIndexChanged: func() {
									currentParty.Type = productTypes[cbType.CurrentIndex()].Code
									saveParty()

								},
							},
							Label{Text: "ПГС1:"},
							newNumberEditPgs(data.Gas1),
							Label{Text: "ПГС2:"},
							newNumberEditPgs(data.Gas2),
							Label{Text: "ПГС3:"},
							newNumberEditPgs(data.Gas3),
							Label{Text: "ПГС4:"},
							newNumberEditPgs(data.Gas4),
						},
					},
				},
			},
		},
	}).Create(); err != nil {
		return err
	}

	pbCancelWork.SetVisible(false)
	mainWindow.Run()

	if err := settings.Save(); err != nil {
		return err
	}
	return nil
}

func runProductDialog(owner walk.Form, p *data.Product) (int, error) {
	var (
		edAddr, edSerial   *walk.NumberEdit
		acceptPB, cancelPB *walk.PushButton
		dlg                *walk.Dialog
	)
	return Dialog{
		Title:      fmt.Sprintf("ДАФ %d", p.ProductID),
		Font:       Font{PointSize: 12, Family: "Segoe UI"},
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		MinSize:    Size{305, 265},
		MaxSize:    Size{305, 265},
		Layout:     VBox{},
		AssignTo:   &dlg,
		Children: []Widget{
			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{},
				Children: []Widget{
					Label{Text: "Адрес:"},
					NumberEdit{
						AssignTo: &edAddr,
						Value:    float64(p.Addr),
						MinValue: 1,
						MaxValue: 127,
						Decimals: 0,
						OnValueChanged: func() {
							p.Addr = modbus.Addr(edAddr.Value())
						},
					},
					Label{Text: "Серийный номер:"},
					NumberEdit{
						AssignTo: &edSerial,
						Value:    float64(p.Serial),
						MinValue: 1,
						MaxValue: 127,
						Decimals: 0,
						OnValueChanged: func() {
							p.Serial = int64(edSerial.Value())
						},
					},

					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{
								Text:     "Ок",
								AssignTo: &acceptPB,
								OnClicked: func() {
									dlg.Accept()
								},
							},
							PushButton{
								Text:     "Отмена",
								AssignTo: &cancelPB,
								OnClicked: func() {
									dlg.Cancel()
								},
							},
						},
					},
				},
			},
		},
	}.Run(owner)
}

var productTypes = []struct {
	Name string
	Code int
}{
	{"ДАФ-М-01", 1},
	{"ДАФ-М-05X", 6},
	{"ДАФ-М-06TPX", 9},
	{"ДАФ-М-08X", 80},
	{"ДАФ-М-08TPX", 85},
}

func indexOfProductTypeCode(productTypeCode int) int {
	for i, x := range productTypes {
		if x.Code == productTypeCode {
			return i
		}
	}
	return -1
}

var mainWindow *walk.MainWindow
