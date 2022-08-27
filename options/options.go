package options

import (
	"fmt"
	"strings"
	"time"

	u "github.com/FT-Labs/physettings/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/FT-Labs/tview"
)

var app *tview.Application
var pages *tview.Pages
var confirm *tview.Modal
var scriptInfo *tview.TextView

var lastFocus tview.Primitive
var lastFocusIndex int = 0
var o1, o2 *tview.Form

func buttonSelGrubTheme() {
    err := u.RunScript(u.POS_GRUB_CHOOSE_THEME)
    if err != nil {
        confirm.SetText(err.Error()).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Succesfully changed grub theme").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorLightGreen)
    }
    lastFocusIndex = 1
    pages.ShowPage("confirm")
}

func buttonSelSddmTheme() {
    err := u.RunScript(u.POS_SDDM_CHOOSE_THEME)
    if err != nil {
        confirm.SetText(err.Error()).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Succesfully changed sddm theme").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorLightGreen)
    }
    lastFocusIndex = 2
    pages.ShowPage("confirm")
}

func buttonSelMakeBar() {
    err := u.RunScript(u.POS_MAKE_BAR)
    if err != nil {
        confirm.SetText(err.Error()).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Succesfully updated statusbar").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorLightGreen)
    }
    lastFocusIndex = 3
    pages.ShowPage("confirm")
}

func checkSelShutdownConfirm(checked bool) {
    if checked {
        u.SetAttribute(u.POWERMENU_CONFIRM, "true")
    } else {
        u.SetAttribute(u.POWERMENU_CONFIRM, "false")
    }
}

func dropSelRofiColor(selection string, i int) {
    u.SetRofiColor(selection)
    confirm.SetText("Rofi colorscheme changed to: " + selection).
            SetBackgroundColor(tcell.Color59).
            SetTextColor(tcell.ColorLightGreen)
    lastFocusIndex = i
    pages.ShowPage("confirm")
}

func dropSelPowerMenuType(selection string, i int) {
    err := u.SetAttribute(u.POWERMENU_TYPE, selection)
    if err != nil {
        confirm.SetText("Failed to set powermenu type").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Powermenu type changed to: " + selection).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorLightGreen)
    }
    lastFocusIndex = i
    pages.ShowPage("confirm")
}

func dropSelPowerMenuStyle(selection string, i int) {
    err := u.SetAttribute(u.POWERMENU_STYLE, selection)
    if err != nil {
        confirm.SetText("Failed to set powermenu style").
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Powermenu style changed to: " + selection).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorLightGreen)
    }
    lastFocusIndex = i
    pages.ShowPage("confirm")
}


func makeDropdown(opt string) *tview.DropDown {
    if opt == u.ROFI_COLOR {
        d := tview.NewDropDown().
                SetLabel("POWERMENU_COLOR : ").
                SetOptions(u.RofiColors, dropSelRofiColor).
                SetCurrentOption(0)
        d.SetFocusFunc(func(){
            go printScriptInfo("Set colorscheme of powermenu.", d)
        })
        return d
    } else if opt == u.POWERMENU_STYLE {
        d := tview.NewDropDown().
                SetLabel(u.POWERMENU_STYLE + " : ").
                SetOptions(u.PowerMenuStyles, dropSelPowerMenuStyle).
                SetCurrentOption(0)
        d.SetFocusFunc(func() {
            go printScriptInfo("Change powermenu style, this will only rearrange items. Look will be similar, but properties will be changed according to powermenu type", d)
        })
        return d
    }// else if opt == POWERMENU_TYPE {
    d := tview.NewDropDown().
            SetLabel(u.POWERMENU_TYPE + " : ").
            SetOptions(u.PowerMenuTypes, dropSelPowerMenuType).
            SetCurrentOption(0)
    d.SetFocusFunc(func() {
        go printScriptInfo("Change type of powermenu. This will change the initial look of powermenu.", d)
    })
    return d
}

func makeScriptsInfoTextView() {
    scriptInfo = tview.NewTextView().
        SetDynamicColors(true).
        SetWordWrap(true).
        SetRegions(true).
        SetChangedFunc(func() {
            app.Draw()
        })
}

func printScriptInfo(s string, p tview.Primitive) {
    scriptInfo.Clear()
    arr := strings.Split(s, " ")
    for _, word := range arr {
        if p.HasFocus() != true {
            scriptInfo.Clear()
            return
        }
        time.Sleep(time.Millisecond * 20)
        fmt.Fprintf(scriptInfo, "%s ", word)
    }
}

// Normal tview requires all of these, not necessary anymore with FT-Labs fork
// type Button struct {
//     *tview.Button
// }

// func (b *Button) GetFieldWidth() int {
//     return len(b.GetLabel())
// }

// func (b *Button) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
//     b.SetExitFunc(handler)
//     return b
// }

// func (b *Button) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
//     b.SetLabelColor(fieldTextColor)
//     b.SetBackgroundColor(fieldBgColor)
//     return b
// }

func makeOptionsForm() *tview.Form {
    c := tview.NewCheckbox().
            SetLabel("ASK ON SHUTDOWN :").
            SetChecked(u.Attrs[u.POWERMENU_CONFIRM] == "true").
            SetChangedFunc(checkSelShutdownConfirm)

    c.SetFocusFunc(func(){
        go printScriptInfo("If selected, when shutting down or rebooting computer, it will ask for a confirmation prompt", c)
    })

    return tview.NewForm().
                SetFieldBackgroundColor(tcell.Color238).
                SetFieldTextColor(tcell.Color255).
                SetLabelColor(tcell.Color33).
                SetItemPadding(3).
                AddDropDown(makeDropdown(u.ROFI_COLOR)).
                AddDropDown(makeDropdown(u.POWERMENU_TYPE)).
                AddDropDown(makeDropdown(u.POWERMENU_STYLE)).
                AddCheckbox(c)
}

func makeScriptsForm() *tview.Form {
    bGrub := tview.NewButton(u.POS_GRUB_CHOOSE_THEME).
                    SetSelectedFunc(buttonSelGrubTheme).
                    SetLabelColorActivated(tcell.Color238)
    bSddm := tview.NewButton(u.POS_SDDM_CHOOSE_THEME).
                    SetSelectedFunc(buttonSelSddmTheme).
                    SetLabelColorActivated(tcell.Color238)
    bBar := tview.NewButton(u.POS_MAKE_BAR).
                    SetSelectedFunc(buttonSelMakeBar).
                    SetLabelColorActivated(tcell.Color238)

    bGrub.SetFocusFunc(func(){
        go printScriptInfo(u.ScriptInfo[u.POS_GRUB_CHOOSE_THEME], bGrub)
    })
    bSddm.SetFocusFunc(func(){
        go printScriptInfo(u.ScriptInfo[u.POS_SDDM_CHOOSE_THEME], bSddm)
    })
    bBar.SetFocusFunc(func(){
        go printScriptInfo(u.ScriptInfo[u.POS_MAKE_BAR], bBar)
    })
    return tview.NewForm().
               SetItemPadding(3).
               SetFieldBackgroundColor(tcell.Color238).
               SetFieldTextColor(tcell.Color255).
               AddButtonItem(bGrub).
               AddButtonItem(bSddm).
               AddButtonItem(bBar)
}


func Options(a *tview.Application,nextSlide func()) (title string, content tview.Primitive){
    app = a
    makeScriptsInfoTextView()
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
                if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
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
                if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
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
                       AddItem(scriptInfo, 0, 3, false).
                       AddItem(tview.NewFlex().
                           SetDirection(tview.FlexColumn).
                           AddItem(o1, 0, 3, true).
                           AddItem(o2, 0, 3, true), 0, 6, true), 0, 0, 1, 1, 0, 0, true)
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
		AddItem(newPrimitive("Use Tab-Shift+Tab or Up-Down keys to navigate, Left-Right to navigate between columns"), 0, 1, false).
		AddItem(newPrimitive("Enter to select (type to search, or use arrow keys), Esc to cancel selection"), 0, 1, false)

	pages.AddPage("flex", flex, true, true).
		AddPage("confirm", confirm, true, false)

    return " 煉 OPTIONS ", pages
}
