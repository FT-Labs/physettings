package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var pages tview.Pages

func selected(selection string, _ int) {
    SetRofiColor(selection)
    pages.ShowPage("modal")
}

func Options(nextSlide func()) (title string, content tview.Primitive){


    pages = *tview.NewPages()
	newPrimitive := func(text string) tview.Primitive {
        if text != "" {
            return tview.NewFrame(nil).
                SetBorders(0, 0, 0, 0, 0, 0).
                AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
        } else {
            dropdown := tview.NewDropDown().
            SetLabel("Select an option (hit Enter): ").
            SetOptions(RofiColors, selected)
            dropdown.SetCurrentOption(0)
            dropdown.SetTitleAlign(tview.AlignCenter)

            flex := tview.NewGrid().SetBorders(true).AddItem(tview.NewFlex().
                SetDirection(tview.FlexRow).
                AddItem(tview.NewBox(), 0, 1, false).
                AddItem(tview.NewFlex().
                    SetDirection(tview.FlexColumn).
                    AddItem(tview.NewBox(), 0, 1, false).
                    AddItem(dropdown, dropdown.GetFieldWidth() + len(dropdown.GetLabel()), dropdown.GetFieldWidth() + len(dropdown.GetLabel()), true).
                    AddItem(tview.NewBox(), 0, 1, false), 0, 2, true).
                AddItem(tview.NewBox(), 0, 1, false), 0, 0, 1, 1, 0, 0, true)
            return flex
        }
	}


	grid := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(tview.NewBox(), 0, 2, false).
        AddItem(newPrimitive("[::b]SET OPTIONS 漣"), 0, 2, false).
        AddItem(tview.NewFlex().
            SetDirection(tview.FlexColumn).
            AddItem(tview.NewBox(), 0, 3, false).
            AddItem(newPrimitive(""), 0, 9, true).
            AddItem(tview.NewBox(), 0, 3, false), 0, 16, true).
		AddItem(newPrimitive("Enter to select (type to search, or use arrow keys)"), 0, 2, false)

	modal := tview.NewModal().
		SetText("Resize the window to see how the grid layout adapts").
		AddButtons([]string{"OK"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		pages.HidePage("modal")
	})


	pages.AddPage("grid", grid, true, true).
		AddPage("modal", modal, false, false)


    return "OPTIONS", grid
}
