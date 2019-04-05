package main

import (
	"context"
	"fmt"
	"github.com/ansel1/merry"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/hako/durafmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/sirupsen/logrus"
	"log"
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
		tblViewProducts *walk.TableView

		neCmd, neArg         *walk.NumberEdit
		pbCancelWork         *walk.PushButton
		lblWork, lblWorkTime *walk.Label

		sbRun      *walk.SplitButton
		gbCmd      *walk.GroupBox
		mainWindow *walk.MainWindow
	)

	lastPartyProductsModel.synchronize = func(f func()) {
		tblViewProducts.Synchronize(f)
	}

	showErr := func(title, text string) {
		walk.MsgBox(mainWindow, title,
			text, walk.MsgBoxIconError|walk.MsgBoxOK)
	}

	delayProgress := new(delayProgressHelp)

	guiShowDelay = delayProgress.Show
	guiHideDelay = delayProgress.Hide

	var workStarted bool
	doWork := func(what string, work func() error) {
		if workStarted {
			panic("already started")
		}
		workStarted = true
		comportContext, cancelComport = context.WithCancel(context.Background())

		sbRun.SetVisible(false)
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

				sbRun.SetVisible(true)
				gbCmd.SetVisible(true)

				pbCancelWork.SetVisible(false)
				lastPartyProductsModel.setInterrogatePlace(-1)
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

	executeProductDialog := func() {
		n := tblViewProducts.CurrentIndex()
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

	productColumns := make([]TableViewColumn, pcConnection+1)

	{
		x := productColumns
		type t = TableViewColumn
		x[pcAddr] = t{Title: "Адрес", Width: 80}
		x[pcSerialNumber] = t{Title: "Номер", Width: 80}
		x[pcProductID] = t{Title: "ID", Width: 80}
		x[pcConcentration] = t{Title: "Концентрация", Width: 150, Precision: 3}
		x[pcCurrent] = t{Title: "Ток", Width: 100, Precision: 1}
		x[pcThreshold1] = t{Title: "Порог 1", Width: 120, Precision: 1}
		x[pcThreshold2] = t{Title: "Порог 2", Width: 120, Precision: 1}
		x[pcMode] = t{Title: "Режим"}
		x[pcFailure] = t{Title: "Отказ"}
		x[pcVersion] = t{Title: "Версия"}
		x[pcGas] = t{Title: "Газ"}
		x[pcConnection] = t{Title: "Связь"}
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

	if err := (MainWindow{
		AssignTo:   &mainWindow,
		Title:      "ЭН8800-6408 Партия ДАФ-М " + data.LastParty().String(),
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

					SplitButton{
						Text:      "Управление",
						AssignTo:  &sbRun,
						MenuItems: menuWorks,
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
						Model:                    lastPartyProductsModel,

						Columns: productColumns,

						ContextMenuItems: []MenuItem{
							Action{
								Text: "Создать новую партию",
							},
							Action{
								Text: "Параметры партии",
								OnTriggered: func() {
									if err := runPartyDialog(mainWindow); err != nil {
										panic(err)
									}
								},
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
							GroupBox{
								AssignTo: &gbCmd,
								Layout:   VBox{},
								Title:    "Команда:",
								Children: []Widget{
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

type delayProgressHelp struct {
	*walk.Composite
	pb     *walk.ProgressBar
	lbl    *walk.Label
	ticker *time.Ticker
	done   chan struct{}
}

func (x *delayProgressHelp) Show(what string, total time.Duration) {
	startMoment := time.Now()
	x.done = make(chan struct{}, 1)
	x.ticker = time.NewTicker(time.Millisecond * 500)
	x.Composite.Synchronize(func() {
		x.SetVisible(true)
		x.pb.SetRange(0, int(total.Nanoseconds()/1000000))
		x.pb.SetValue(0)
		_ = x.lbl.SetText(fmt.Sprintf("%s: %s", what, durafmt.Parse(total)))
	})
	go func() {
		defer func() {
			logrus.Debugln("timer closed")
		}()
		for {
			select {
			case <-x.ticker.C:
				x.Composite.Synchronize(func() {
					x.pb.SetValue(int(time.Since(startMoment).Nanoseconds() / 1000000))
				})
			case <-x.done:
				return
			}
		}
	}()
}

func (x *delayProgressHelp) Hide() {
	x.ticker.Stop()
	close(x.done)
	x.Composite.Synchronize(func() {
		x.SetVisible(false)
	})

}

func (x *delayProgressHelp) Widget() Widget {
	return Composite{
		AssignTo: &x.Composite,
		Layout:   HBox{},
		Visible:  false,
		Children: []Widget{
			Label{AssignTo: &x.lbl},
			ScrollView{
				Layout:        VBox{SpacingZero: true, MarginsZero: true},
				VerticalFixed: true,
				Children: []Widget{
					ProgressBar{
						AssignTo: &x.pb,
						MaxSize:  Size{0, 15},
						MinSize:  Size{0, 15},
					},
				},
			},

			PushButton{
				Text: "Продолжить без задержки",
				OnClicked: func() {
					skipDelay()
				},
			},
		},
	}
}
