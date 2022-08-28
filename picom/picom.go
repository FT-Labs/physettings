package picom

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

func checkSelAnimConfirm(checked bool) {
    if checked {
        changePicomAttribute(_animations, "true")
    } else {
        changePicomAttribute(_animations, "false")
    }
}

func checkSelExperimental(checked bool) {
    if checked {
        u.SetAttribute(u.PICOM_EXPERIMENTAL, "true")
        changePicomAttribute(u.PICOM_EXPERIMENTAL, "true")
    } else {
        u.SetAttribute(u.PICOM_EXPERIMENTAL, "false")
        changePicomAttribute(u.PICOM_EXPERIMENTAL, "false")
    }
}

func checkSelFadeConfirm(checked bool) {
    if checked {
        changePicomAttribute(_fading, "true")
    } else {
        changePicomAttribute(_fading, "false")
    }
}

func checkSelFadeNextTag(checked bool) {
    if checked {
        changePicomAttribute(_enable_fading_next_tag, "true")
    } else {
        changePicomAttribute(_enable_fading_next_tag, "false")
    }
}

func checkSelFadePrevTag(checked bool) {
    if checked {
        changePicomAttribute(_enable_fading_prev_tag, "true")
    } else {
        changePicomAttribute(_enable_fading_prev_tag, "false")
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

func makeOptionsForm() *tview.Form {

    c := tview.NewCheckbox().
            SetLabel("EXPERIMENTAL BACKENDS:").
            SetChecked(u.Attrs[u.PICOM_EXPERIMENTAL] == "true").
            SetChangedFunc(checkSelExperimental)

    c.SetFocusFunc(func(){
        go printScriptInfo("Enable experimental backends in picom.\nThis will add dual-kawase blur, which is preferred.\nNote that this uses GLX backend, which doesn't work correctly in legacy hardware.", c)
    })

    c1 := tview.NewCheckbox().
            SetLabel("ENABLE ANIMATIONS :").
            SetChecked(picomOpts[_animations] == "true").
            SetChangedFunc(checkSelAnimConfirm)

    c1.SetFocusFunc(func(){
        go printScriptInfo("Enable or disable animations in picom", c1)
    })

    c2 := tview.NewCheckbox().
            SetLabel("ENABLE FADING :").
            SetChecked(picomOpts[_fading] == "true").
            SetChangedFunc(checkSelFadeConfirm)

    c2.SetFocusFunc(func(){
        go printScriptInfo("Enable fading when opening - closing windows.\nWindows will go from transparent to opaque if set.", c2)
    })

    c3 := tview.NewCheckbox().
            SetLabel("NEXT TAG FADING :").
            SetChecked(picomOpts[_enable_fading_next_tag] == "true").
            SetChangedFunc(checkSelFadeNextTag)

    c3.SetFocusFunc(func(){
        go printScriptInfo("Enable fading for incoming tag.\nNew windows that are coming from next tag will go from transparent to opaque.", c3)
    })

    c4 := tview.NewCheckbox().
            SetLabel("PREV TAG FADING :").
            SetChecked(picomOpts[_enable_fading_prev_tag] == "true").
            SetChangedFunc(checkSelFadePrevTag)

    c4.SetFocusFunc(func(){
        go printScriptInfo("Enable fading for next tag.\nOld windows that are going out from current tag will go from opaque to transparent.", c4)
    })

    i1 := tview.NewInputField().
            SetLabel("ANIM SPEED IN TAG :").
            SetAcceptanceFunc(tview.InputFieldFloat).
            SetPlaceholder(picomOpts[_animation_stiffness_in_tag]).
            SetFieldWidth(3)

    i1.SetFocusFunc(func(){
        go printScriptInfo("Set animation speed in current tag.\n125 is default. Enter an integer or float number.", c4)
    })


    return tview.NewForm().
                SetFieldBackgroundColor(tcell.Color238).
                SetFieldTextColor(tcell.Color255).
                SetLabelColor(tcell.Color33).
                SetItemPadding(2).
                AddCheckbox(c).
                AddCheckbox(c1).
                AddCheckbox(c2).
                AddCheckbox(c3).
                AddCheckbox(c4).
                AddInputFieldItem(i1)
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


func Picom(a *tview.Application,nextSlide func()) (title string, content tview.Primitive){
    readPicomOpts()
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

    return " 𧻓 PICOM - ANIMATIONS ", pages
}
