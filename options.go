package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var pages *tview.Pages
var confirm *tview.Modal
var lastFocus tview.Primitive
var lastFocusIndex int
var o1, o2 *tview.Form

func dropSelRofiColor(selection string, i int) {
    SetRofiColor(selection)
    confirm.SetText("Rofi colorscheme changed to: " + selection).
            SetBackgroundColor(tcell.Color59)
    lastFocusIndex = i
    pages.ShowPage("confirm")
}

func dropSelPowerMenuType(selection string, i int) {
    err := SetAttribute(POWERMENU_TYPE, selection)
    if err != nil {
        confirm.SetText("Failed to set powermenu type").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Powermenu type changed to: " + selection).
                SetBackgroundColor(tcell.Color59)
    }
    lastFocusIndex = i
    pages.ShowPage("confirm")
}

func dropSelPowerMenuStyle(selection string, i int) {
    err := SetAttribute(POWERMENU_STYLE, selection)
    if err != nil {
        confirm.SetText("Failed to set powermenu style").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Powermenu style changed to: " + selection).
                SetBackgroundColor(tcell.Color59)
    }
    lastFocusIndex = i
    pages.ShowPage("confirm")
}

func buttonSelGrubTheme() {
    const c = "pOS-grub-choose-theme"
    err := RunScript(c)
    if err != nil {
        confirm.SetText(err.Error()).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Succesfully changed grub theme").
                SetBackgroundColor(tcell.Color59)
    }
    lastFocusIndex = 1
    pages.ShowPage("confirm")
}

func buttonSelSddmTheme() {
    const c = "pOS-sddm-choose-theme"
    err := RunScript(c)
    if err != nil {
        confirm.SetText(err.Error()).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Succesfully changed sddm theme").
                SetBackgroundColor(tcell.Color59)
    }
    lastFocusIndex = 2
    pages.ShowPage("confirm")
}

func buttonSelMakeBar() {
    const c = "pOS-make-bar"
    err := RunScript(c)
    if err != nil {
        confirm.SetText(err.Error()).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Succesfully updated statusbar.").
                SetBackgroundColor(tcell.Color59)
    }
    lastFocusIndex = 3
    pages.ShowPage("confirm")
}

func makeDropdown(opt string) (l string, o []string, idx int, sel func(o string, oindex int)) {
    if opt == ROFI_COLOR {
        return ROFI_COLOR + " : ",
                    RofiColors,
                    0,
                    dropSelRofiColor
    } else if opt == POWERMENU_STYLE {
        return       POWERMENU_STYLE + " : ",
                     PowerMenuStyles,
                     0,
                     dropSelPowerMenuStyle
    }// else if opt == POWERMENU_TYPE {
    return POWERMENU_TYPE + " : ",
           PowerMenuTypes,
           0,
           dropSelPowerMenuType
}

type Button struct {
    *tview.Button
}

func (b Button) GetFieldWidth() int {
    return 0
}

func (b Button) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
    b.SetExitFunc(handler)
    return b
}

func (b Button) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
    b.SetLabelColor(fieldTextColor)
    b.SetBackgroundColor(fieldBgColor)
    return b
}


func makeOptionsForm() *tview.Form {
    return tview.NewForm().
                SetFieldBackgroundColor(tcell.Color16).
                SetFieldTextColor(tcell.Color231).
                SetItemPadding(3).
                AddDropDown(makeDropdown(ROFI_COLOR)).
                AddDropDown(makeDropdown(POWERMENU_TYPE)).
                AddDropDown(makeDropdown(POWERMENU_STYLE))
}

func makeScriptsForm() *tview.Form {
    return tview.NewForm().
               SetItemPadding(3).
               SetFieldBackgroundColor(tcell.Color16).
               SetFieldTextColor(tcell.Color231).
               AddFormItem(Button{tview.NewButton("SET GRUB THEME").SetSelectedFunc(buttonSelGrubTheme)}).
               AddFormItem(Button{tview.NewButton("SET SDDM THEME").SetSelectedFunc(buttonSelSddmTheme)}).
               AddFormItem(Button{tview.NewButton("MAKE STATUSBAR").SetSelectedFunc(buttonSelMakeBar)})
}


func Options(nextSlide func()) (title string, content tview.Primitive){

	confirm = tview.NewModal().
		AddButtons([]string{"OK"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		pages.HidePage("confirm")
        if lastFocus != nil && lastFocus != app.GetFocus() {
            app.SetFocus(lastFocus)
        }
	})

    pages = tview.NewPages()
	newPrimitive := func(text string) tview.Primitive {
        if text != "" {
            return tview.NewFrame(nil).
                SetBorders(0, 0, 0, 0, 0, 0).
                AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
        } else {
            o1 = makeOptionsForm()
            o2 = makeScriptsForm()

            o1.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
                if event.Key() == tcell.KeyBacktab {
                    app.SetFocus(o2)
                    lastFocus = o2
                    return nil
                }
                return event
            })

            o1.SetFocusFunc(func() {
                o1.SetFocus(lastFocusIndex)
            })

            o2.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
                if event.Key() == tcell.KeyBacktab {
                    app.SetFocus(o1)
                    lastFocus = o1
                    return nil
                }
                return event
            })

            o2.SetFocusFunc(func() {
                o2.SetFocus(lastFocusIndex)
            })


            return tview.NewGrid().
                       SetBordersColor(tcell.Color33).
                       SetBorders(true).
                       AddItem(tview.NewFlex().
                       SetDirection(tview.FlexRow).
                       AddItem(tview.NewBox(), 0, 1, false).
                       AddItem(tview.NewFlex().
                           SetDirection(tview.FlexColumn).
                           AddItem(o1, 0, 3, true).
                           AddItem(o2, 0, 3, true), 0, 6, true).
                       AddItem(tview.NewBox(), 0, 1, false), 0, 0, 1, 1, 0, 0, true)
        }
	}

	flex := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(tview.NewBox(), 0, 1, false).
        AddItem(newPrimitive("[::b]SET OPTIONS 漣"), 0, 1, false).
        AddItem(tview.NewFlex().
            SetDirection(tview.FlexColumn).
            AddItem(tview.NewBox(), 0, 3, false).
            AddItem(newPrimitive(""), 0, 9, true).
            AddItem(tview.NewBox(), 0, 3, false), 0, 16, true).
		AddItem(newPrimitive("Press Tab to navigate in current column, Shift+Tab to switch between columns"), 0, 1, false).
		AddItem(newPrimitive("Enter to select (type to search, or use arrow keys)"), 0, 1, false)

	pages.AddPage("flex", flex, true, true).
		AddPage("confirm", confirm, true, false)

    return "OPTIONS", pages
}
