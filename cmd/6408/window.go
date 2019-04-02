package main

import (
	"fmt"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
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
		mainWindow          *walk.MainWindow
		comboBoxProductType,
		comboBoxComportProducts,
		comboBoxComportHart *walk.ComboBox

		numberEditPgs [4]*walk.NumberEdit
	)

	var currentParty = data.LastParty()

	saveParty := func() {
		if err := data.SaveParty(currentParty); err != nil {
			walk.MsgBox(mainWindow, "Ошибка данных", fmt.Sprintf("%v: %v", currentParty, err), walk.MsgBoxIconError|walk.MsgBoxOK)
		}
	}

	newNumberEditPgs := func(n data.Gas) NumberEdit {
		return NumberEdit{
			Value:    currentParty.Pgs(n),
			AssignTo: &numberEditPgs[n-1],
			MinValue: 0,
			Decimals: 2,
			OnValueChanged: func() {
				currentParty.SetPgs(n, numberEditPgs[n-1].Value())
				saveParty()
			},
		}
	}

	if err := (MainWindow{
		AssignTo:   &mainWindow,
		Title:      "ЭН8800-6408",
		Name:       "MainWindow",
		Font:       Font{PointSize: 14, Family: "Segoe UI"},
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Size:       Size{800, 600},
		Layout:     VBox{},

		Children: []Widget{
			ScrollView{
				VerticalFixed: true,
				Layout:        HBox{},
				Children: []Widget{
					GroupBox{
						Title:  "-",
						Layout: HBox{},
					},
				},
			},
			ScrollView{
				Layout: HBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					TableView{
						Font:                     Font{PointSize: 14, Family: "Segoe UI"},
						NotSortableByHeaderClick: true,
						CheckBoxes:               true,
						Columns: []TableViewColumn{
							{Title: "Адрес", Width: 80},
							{Title: "Заводской номер", Width: 150},
							{Title: "Концентрация", Width: 150},
							{Title: "Ток", Width: 100},
							{Title: "Порог 1", Width: 120},
							{Title: "Порог 2", Width: 120},
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
									newComboBoxComport(&comboBoxComportProducts, "comport_products"),
									Label{Text: "HART модем:"},
									newComboBoxComport(&comboBoxComportHart, "comport_hart"),
								},
							},
							Label{Text: "Исполнение:"},
							ComboBox{
								Model:         productTypes,
								AssignTo:      &comboBoxProductType,
								DisplayMember: "Name",
								CurrentIndex:  indexOfProductTypeCode(currentParty.Type),
								OnCurrentIndexChanged: func() {
									currentParty.Type = productTypes[comboBoxProductType.CurrentIndex()].Code
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
	mainWindow.Run()
	if err := settings.Save(); err != nil {
		return err
	}
	return nil
}

var productTypes = []struct {
	Name string;
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
