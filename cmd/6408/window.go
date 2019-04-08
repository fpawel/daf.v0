package main

import (
	"context"
	"fmt"
	"github.com/ansel1/merry"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/daf/internal/viewmodel"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"math"
	"time"
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
		cbComportDaf,
		cbComportHart *walk.ComboBox
		tblViewProducts,
		tblViewProductValues,
		tblViewProductEntries *walk.TableView

		neCmd, neArg         *walk.NumberEdit
		pbCancelWork         *walk.PushButton
		lblWork, lblWorkTime *walk.Label
		btnRun               *walk.SplitButton
		gbCmd                *walk.GroupBox
		mainWindow           *walk.MainWindow
	)

	prodsMdl = viewmodel.NewDafProductsTable(func(f func()) {
		tblViewProducts.Synchronize(f)
	})

	showErr := func(title, text string) {
		walk.MsgBox(mainWindow, title,
			text, walk.MsgBoxIconError|walk.MsgBoxOK)
	}

	delayProgress := new(delayHelp)

	guiShowDelay = delayProgress.Show
	guiHideDelay = delayProgress.Hide

	var workStarted bool
	doWork := func(what string, work func() error) {
		if workStarted {
			panic("already started")
		}
		workStarted = true

		prodsMdl.ClearConnectionsInfo()

		comportContext, cancelComport = context.WithCancel(context.Background())
		btnRun.SetVisible(false)
		gbCmd.SetVisible(false)

		portDaf.PortName = cbComportDaf.Text()
		portHart.PortName = cbComportHart.Text()
		pbCancelWork.SetVisible(true)
		_ = lblWorkTime.SetText(time.Now().Format("15:04:05"))
		_ = lblWork.SetText(fmt.Sprintf("%s: выполняется", what))
		lblWork.SetTextColor(walk.RGB(128, 0, 0))

		go func() {
			err := work()

			_ = portHart.Port.Close()
			_ = portDaf.Port.Close()
			mainWindow.Synchronize(func() {
				workStarted = false

				gbCmd.SetVisible(true)
				btnRun.SetVisible(true)

				pbCancelWork.SetVisible(false)
				prodsMdl.SetInterrogatePlace(-1)
				_ = lblWorkTime.SetText(time.Now().Format("15:04:05"))
				if err != nil {
					if merry.Is(err, context.Canceled) {
						lblWork.SetTextColor(walk.RGB(139, 69, 19))
						_ = lblWork.SetText(fmt.Sprintf("%s: прервано", what))
					} else {
						lblWork.SetTextColor(walk.RGB(255, 0, 0))
						_ = lblWork.SetText(fmt.Sprintf("%s: %v", what, err))
						showErr(what, err.Error())
					}

				} else {
					lblWork.SetTextColor(walk.RGB(0, 0, 128))
					_ = lblWork.SetText(fmt.Sprintf("%s: выполнено", what))
				}

			})
		}()
	}

	menuWorks := []MenuItem{
		Action{
			Text: "Опрос",
			OnTriggered: func() {
				doWork("опрос", interrogateProducts)
			},
		},
		Action{
			Text: "Настройка токового выхода",
			OnTriggered: func() {
				doWork("настройка токового выхода", setupCurrents)
			},
		},
		Separator{},
		Action{
			Text: "отключить газ",
			OnTriggered: func() {
				doWork("отключить газ", func() error {
					return switchGas(0)
				})
			},
		},

		Action{
			Text: "задержка",
			OnTriggered: func() {
				doWork("некоторая задержка", func() error {

					if err := delay("sdfsdf", time.Minute); err != nil {
						return err
					}

					return delay("rtyrty", time.Minute)
				})
			},
		},
	}

	for gas := data.Gas1; gas < 5; gas++ {
		gas := gas
		s := fmt.Sprintf("газ %d", gas)
		menuWorks = append(menuWorks, Action{
			Text: s,
			OnTriggered: func() {
				doWork(s, func() error {
					return switchGas(gas)
				})
			},
		})
	}

	prodValuesMdl := viewmodel.NewDafProductValuesTable()

	validateProductValuesTable := func() {
		n := tblViewProducts.CurrentIndex()
		if n > -1 && n < prodsMdl.RowCount() {
			prodValuesMdl.SetProduct(prodsMdl.ProductAt(n).ProductID)
			if prodValuesMdl.RowCount() > 0 {
				tblViewProductValues.SetVisible(true)
				return
			}
		}
		tblViewProductValues.SetVisible(false)
	}

	prodEntriesMdl := viewmodel.NewDafProductEntriesTable()

	validateProductEntriesTable := func() {
		n := tblViewProducts.CurrentIndex()
		if n > -1 && n < prodsMdl.RowCount() {
			prodEntriesMdl.SetProduct(prodsMdl.ProductAt(n).ProductID)
			if prodEntriesMdl.RowCount() > 0 {
				tblViewProductEntries.SetVisible(true)
				return
			}
		}
		tblViewProductEntries.SetVisible(false)
	}

	if err := (MainWindow{
		AssignTo: &mainWindow,
		Title: "ЭН8800-6408 Партия ДАФ-М " + (func() string {
			p := data.GetLastParty()
			return fmt.Sprintf("№%d %s", p.PartyID, p.CreatedAt.Format("02.01.2006"))
		}()),
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
					SplitButton{
						Text: "Партия",
						MenuItems: []MenuItem{
							Action{
								Text: "Создать новую",
								OnTriggered: func() {

								},
							},
							Action{
								Text: "Параметры",
								OnTriggered: func() {
									runPartyDialog(mainWindow)
								},
							},
							Action{
								Text: "Добавить прибор в партию",
								Shortcut: Shortcut{
									Key: walk.KeyInsert,
								},
								OnTriggered: func() {
									prodsMdl.AddNewProduct()
								},
							},
						},
					},
					SplitButton{
						Text:      "Управление",
						AssignTo:  &btnRun,
						MenuItems: menuWorks,
					},
					PushButton{
						AssignTo: &pbCancelWork,
						Text:     "Прервать",
						OnClicked: func() {
							cancelComport()
						},
					},

					Label{
						AssignTo:  &lblWorkTime,
						TextColor: walk.RGB(0, 128, 0),
					},
					Label{
						AssignTo: &lblWork,
					},
					delayProgress.Widget(),
				},
			},
			ScrollView{
				Layout: HBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					TableView{
						AssignTo:                 &tblViewProducts,
						NotSortableByHeaderClick: true,
						LastColumnStretched:      true,
						CheckBoxes:               true,
						Model:                    prodsMdl,
						OnItemActivated: func() {
							n := tblViewProducts.CurrentIndex()
							if n < 0 || n >= prodsMdl.RowCount() {
								return
							}
							runProductDialog(mainWindow, prodsMdl.ProductAt(n))
							prodsMdl.PublishRowChanged(n)
						},
						OnKeyDown: func(key walk.Key) {
							switch key {

							case walk.KeyInsert:
								m := prodsMdl
								m.AddNewProduct()
								runProductDialog(mainWindow, m.ProductAt(m.RowCount()-1))
								prodsMdl.PublishRowChanged(m.RowCount() - 1)

							case walk.KeyDelete:
								n := tblViewProducts.CurrentIndex()
								m := prodsMdl
								if n < 0 || n >= m.RowCount() {
									return
								}
								if err := data.DBProducts.Delete(m.ProductAt(n)); err != nil {
									showErr("Ошибка данных", err.Error())
								}
								prodsMdl.Validate()
							}

						},
						OnCurrentIndexChanged: func() {
							validateProductValuesTable()
							validateProductEntriesTable()
						},

						Columns: viewmodel.ProductColumns,
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
							GroupBox{
								AssignTo: &gbCmd,
								Layout:   VBox{},
								Title:    "Команда:",
								Children: []Widget{
									Label{Text: "Код:"},
									NumberEdit{
										AssignTo: &neCmd,
										MinValue: 1,
										MaxValue: math.MaxFloat64,
									},
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
				},
			},
			Composite{
				Layout: HBox{SpacingZero: true, MarginsZero: true},
				Children: []Widget{
					TableView{
						AssignTo:                 &tblViewProductValues,
						NotSortableByHeaderClick: true,
						LastColumnStretched:      true,
						Model:                    prodValuesMdl,
						Columns:                  viewmodel.ProductValueColumns,
					},
					TableView{
						AssignTo:                 &tblViewProductEntries,
						NotSortableByHeaderClick: true,
						LastColumnStretched:      true,
						Model:                    prodEntriesMdl,
						Columns:                  viewmodel.ProductEntryColumns,
					},
				},
			},
		},
	}).Create(); err != nil {
		return err
	}

	pbCancelWork.SetVisible(false)
	prodsMdl.Validate()
	mainWindow.Run()

	if err := settings.Save(); err != nil {
		return err
	}
	return nil
}

func runProductDialog(owner walk.Form, p *data.Product) {
	var (
		edAddr, edSerial *walk.NumberEdit
		dlg              *walk.Dialog
		btn              *walk.PushButton
		lblError         *walk.Label
		saveOnEdit       = false
	)

	save := func(what string) {
		if !saveOnEdit {
			return
		}
		p.Serial = int64(edSerial.Value())
		p.Addr = modbus.Addr(edAddr.Value())
		_ = lblError.SetText("")
		if err := data.DBProducts.Save(p); err != nil {
			_ = lblError.SetText(fmt.Sprintf("%s: дублирование значения: %v", what, err))
			if err := data.DBProducts.FindByPrimaryKeyTo(p, p.ProductID); err != nil {
				panic(err)
			}
		}
		if edSerial.Value() != float64(p.Serial) {
			edSerial.SetTextColor(0xFF)
		} else {
			edSerial.SetTextColor(0)
		}

		if edAddr.Value() != float64(p.Addr) {
			edAddr.SetTextColor(0xFF)
		} else {
			edAddr.SetTextColor(0)
		}

	}
	d := Dialog{
		Title:         fmt.Sprintf("ДАФ %d", p.ProductID),
		Font:          Font{PointSize: 12, Family: "Segoe UI"},
		Background:    SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Layout:        Grid{Columns: 2},
		AssignTo:      &dlg,
		DefaultButton: &btn,
		CancelButton:  &btn,
		Children: []Widget{
			Label{Text: "Адрес:", TextAlignment: AlignFar},
			NumberEdit{
				AssignTo: &edAddr,
				Value:    float64(p.Addr),
				MinValue: 1,
				MaxValue: 127,
				Decimals: 0,
				OnValueChanged: func() {
					save(fmt.Sprintf("адрес: %v", edAddr.Value()))
				},
			},
			Label{Text: "Серийный номер:", TextAlignment: AlignFar},
			NumberEdit{
				AssignTo: &edSerial,
				Value:    float64(p.Serial),
				MinValue: 1,
				MaxValue: math.MaxFloat64,
				Decimals: 0,
				OnValueChanged: func() {
					save(fmt.Sprintf("cерийный номер: %v", edSerial.Value()))
				},
			},
			Composite{},
			PushButton{
				AssignTo: &btn,
				Text:     "Закрыть",
				OnClicked: func() {
					dlg.Accept()
				},
			},
			Label{
				ColumnSpan: 2,
				AssignTo:   &lblError,
				TextColor:  0xFF,
			},
		},
	}
	if err := d.Create(owner); err != nil {
		walk.MsgBox(owner, fmt.Sprintf("ДАФ %d", p.ProductID), err.Error(), walk.MsgBoxIconError|walk.MsgBoxOK)
		return
	}
	saveOnEdit = true
	dlg.Run()
}
