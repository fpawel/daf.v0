package main

import (
	"github.com/fpawel/daf/internal/data"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runPartyDialog(owner walk.Form) error {
	var (
		dlg           *walk.Dialog
		cbGas, cbType *walk.ComboBox
		btn           *walk.PushButton
		saveEditParty bool
		party         = data.GetLastParty()
	)

	saveParty := func() {
		if !saveEditParty {
			return
		}
		if err := data.DBProducts.Save(party); err != nil {
			walk.MsgBox(dlg, "Ошибка данных",
				err.Error(), walk.MsgBoxIconError|walk.MsgBoxOK)
		}
	}

	type NameCode struct {
		Name string
		Code int
	}

	dafTypes := []NameCode{
		{"ДАФ-М-01", 1},
		{"ДАФ-М-05X", 6},
		{"ДАФ-М-06TPX", 9},
		{"ДАФ-М-08X", 80},
		{"ДАФ-М-08TPX", 85},
	}

	dafGases := []string{
		"ацетон С₃Н₆",
		"гексан C₆H₁₄",
		"бензол С₆Н₆",
		"стирол С₈Н₈",
		"толуол С₆Н₅СН₃",
		"фенол С₆Н₆О",
		"этанол С₂Н₅ОН",
		"циклогексан С₆Н₁₂",
	}

	nePartyField := func(pValue *float64, minValue float64, decimals int) NumberEdit {
		var ne *walk.NumberEdit
		return NumberEdit{
			Value:    *pValue,
			AssignTo: &ne,
			MinValue: minValue,
			Decimals: decimals,
			OnValueChanged: func() {
				*pValue = ne.Value()
				saveParty()
			},
		}
	}

	dialog := Dialog{
		Font:          Font{PointSize: 12, Family: "Segoe UI"},
		AssignTo:      &dlg,
		Layout:        Grid{Columns: 2},
		DefaultButton: &btn,
		CancelButton:  &btn,
		Children: []Widget{
			Label{Text: "Исполнение:", TextAlignment: AlignFar},
			ComboBox{
				Model:         dafTypes,
				AssignTo:      &cbType,
				DisplayMember: "Name",
				CurrentIndex: func() int {
					for i, x := range dafTypes {
						if x.Code == party.Type {
							return i
						}
					}
					return -1
				}(),
				OnCurrentIndexChanged: func() {
					party.Type = dafTypes[cbType.CurrentIndex()].Code
					saveParty()
				},
			},

			Label{Text: "Компонент:", TextAlignment: AlignFar},
			ComboBox{
				Model:    dafGases,
				AssignTo: &cbGas,
				CurrentIndex: func() int {
					for i, x := range dafGases {
						if x == party.Component {
							return i
						}
					}
					return -1
				}(),
				OnCurrentIndexChanged: func() {
					party.Component = dafGases[cbGas.CurrentIndex()]
					saveParty()
				},
			},

			Label{Text: "Диапазон:", TextAlignment: AlignFar},
			nePartyField(&party.Scale, 0, 0),

			Label{Text: "Дапазон абс. погр.:", TextAlignment: AlignFar},
			nePartyField(&party.AbsoluteErrorRange, 0, 0),
			Label{Text: "Предел абс. погр.:", TextAlignment: AlignFar},
			nePartyField(&party.AbsoluteErrorLimit, 0, 0),
			Label{Text: "Предел отн. погр., %:", TextAlignment: AlignFar},
			nePartyField(&party.RelativeErrorLimit, 0, 0),

			Label{Text: "ПГС1:", TextAlignment: AlignFar},
			nePartyField(&party.Pgs1, 0, 0),
			Label{Text: "ПГС2:", TextAlignment: AlignFar},
			nePartyField(&party.Pgs2, 0, 0),
			Label{Text: "ПГС3:", TextAlignment: AlignFar},
			nePartyField(&party.Pgs3, 0, 0),
			Label{Text: "ПГС4:", TextAlignment: AlignFar},
			nePartyField(&party.Pgs4, 0, 0),

			Label{Text: "Порог 1:", TextAlignment: AlignFar},
			nePartyField(&party.Threshold1Production, 0, 0),
			Label{Text: "Порог 2:", TextAlignment: AlignFar},
			nePartyField(&party.Threshold2Production, 0, 0),

			Label{Text: "Порог 1, настройка:", TextAlignment: AlignFar},
			nePartyField(&party.Threshold1Test, 0, 0),
			Label{Text: "Порог 2, настройка:", TextAlignment: AlignFar},
			nePartyField(&party.Threshold2Test, 0, 0),

			Composite{},
			PushButton{
				AssignTo: &btn,
				Text:     "Закрыть",
				OnClicked: func() {
					dlg.Accept()
				},
			},
		},
	}
	if err := dialog.Create(owner); err != nil {
		return err
	}
	saveEditParty = true
	dlg.Run()
	return nil
}
