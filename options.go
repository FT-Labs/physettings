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

            flex := tview.NewFlex().
                AddItem(tview.NewBox(), 0, 1, false).
                AddItem(tview.NewFlex().
                    SetDirection(tview.FlexColumn).
                    AddItem(tview.NewBox(), 0, 1, false).
                    AddItem(dropdown, 0, 2, true).
                    AddItem(tview.NewBox(), 0, 1, false), 0, 2, true).
                AddItem(tview.NewBox(), 0, 1, false)
            return flex
        }
	}


	grid := tview.NewGrid().
		SetRows(-1, -1, -1, -7, -1).
		SetColumns(-1, -8, -1).
        SetBorders(true).
        AddItem(newPrimitive(""), 2, 1, 3, 1, 0, 0, true)

    grid.SetBorders(false)
    grid.AddItem(newPrimitive("[::b]SET OPTIONS 漣"), 1, 0, 1, 3, 0, 0, true).
		AddItem(newPrimitive("Enter to select (type to search, or use arrow keys)"), 5, 0, 1, 3, 0, 0, false)

	modal := tview.NewModal().
		SetText("Resize the window to see how the grid layout adapts").
		AddButtons([]string{"OK"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		pages.HidePage("modal")
	})


	pages.AddPage("grid", grid, true, true).
		AddPage("modal", modal, false, false)


    return "OPTIONS", grid
}
