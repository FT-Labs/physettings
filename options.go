package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var pages *tview.Pages
var confirm *tview.Modal

func selRofiColor(selection string, _ int) {
    SetRofiColor(selection)
    confirm.SetText("Rofi colorscheme changed to: " + selection).
            SetBackgroundColor(tcell.Color59)
    pages.ShowPage("confirm")
}

func selPowerMenuType(selection string, _ int) {
    err := SetAttribute(POWERMENU_TYPE, selection)
    if err != nil {
        confirm.SetText("Failed to set powermenu type").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Powermenu type changed to: " + selection).
                SetBackgroundColor(tcell.Color59)
    }
    pages.ShowPage("confirm")
}

func selPowerMenuStyle(selection string, _ int) {
    err := SetAttribute(POWERMENU_STYLE, selection)
    if err != nil {
        confirm.SetText("Failed to set powermenu style").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Powermenu style changed to: " + selection).
                SetBackgroundColor(tcell.Color59)
    }
    pages.ShowPage("confirm")
}

func makeDropdown(opt string) (l string, o []string, idx int, sel func(o string, oindex int)) {
    if opt == ROFI_COLOR {
        return ROFI_COLOR + ":    ",
                    RofiColors,
                    0,
                    selRofiColor
    } else if opt == POWERMENU_STYLE {
        return       POWERMENU_STYLE + ":    ",
                     PowerMenuStyles,
                     0,
                     selPowerMenuStyle
    }// else if opt == POWERMENU_TYPE {
        return       POWERMENU_TYPE + "    :",
                     PowerMenuTypes,
                     0,
                     selPowerMenuType
}

func makeOptionsForm() *tview.Form {
    f := tview.NewForm().
        AddDropDown(makeDropdown(ROFI_COLOR)).
        AddDropDown(makeDropdown(POWERMENU_TYPE)).
        AddDropDown(makeDropdown(POWERMENU_STYLE))
    return f
}


func Options(nextSlide func()) (title string, content tview.Primitive){

	confirm = tview.NewModal().
		AddButtons([]string{"OK"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		pages.HidePage("confirm")
	})

    pages = tview.NewPages()
	newPrimitive := func(text string) tview.Primitive {
        if text != "" {
            return tview.NewFrame(nil).
                SetBorders(0, 0, 0, 0, 0, 0).
                AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
        } else {

            grid := tview.NewGrid().
                SetBordersColor(tcell.Color33).
                SetBorders(true).
                AddItem(tview.NewFlex().
                SetDirection(tview.FlexRow).
                AddItem(tview.NewBox(), 0, 1, false).
                AddItem(tview.NewFlex().
                    SetDirection(tview.FlexColumn).
                    AddItem(makeOptionsForm(), 0, 6, true).
                //     AddItem(tview.NewBox(), 0, 1, false).
                //     AddItem(&dropSlice[0],
                //                  dropSlice[0].GetFieldWidth() + len(dropSlice[0].GetLabel()),
                //                  dropSlice[0].GetFieldWidth() + len(dropSlice[0].GetLabel()), true).
                //     AddItem(tview.NewBox(), 0, 1, false), 0, 4, true).
                // AddItem(tview.NewFlex().
                //     SetDirection(tview.FlexColumn).
                //     AddItem(tview.NewBox(), 0, 1, false).
                //     AddItem(&dropSlice[1],
                //                  dropSlice[0].GetFieldWidth() + len(dropSlice[0].GetLabel()),
                //                  dropSlice[0].GetFieldWidth() + len(dropSlice[0].GetLabel()), true).
                //     AddItem(tview.NewBox(), 0, 1, false), 0, 4, true).
                // AddItem(tview.NewFlex().
                //     SetDirection(tview.FlexColumn).
                //     AddItem(tview.NewBox(), 0, 1, false).
                //     AddItem(&dropSlice[2],
                //                  dropSlice[0].GetFieldWidth() + len(dropSlice[0].GetLabel()),
                //                  dropSlice[0].GetFieldWidth() + len(dropSlice[0].GetLabel()), true).
                    AddItem(tview.NewBox(), 0, 1, false), 0, 4, true).
                AddItem(tview.NewBox(), 0, 1, false), 0, 0, 1, 1, 0, 0, true)

            return grid
        }
	}

	flex := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(tview.NewBox(), 0, 2, false).
        AddItem(newPrimitive("[::b]SET OPTIONS 漣"), 0, 2, false).
        AddItem(tview.NewFlex().
            SetDirection(tview.FlexColumn).
            AddItem(tview.NewBox(), 0, 3, false).
            AddItem(newPrimitive(""), 0, 9, true).
            AddItem(tview.NewBox(), 0, 3, false), 0, 16, true).
		AddItem(newPrimitive("Enter to select (type to search, or use arrow keys)"), 0, 2, false)

	pages.AddPage("flex", flex, true, true).
		AddPage("confirm", confirm, false, false)

    return "OPTIONS", pages
}
