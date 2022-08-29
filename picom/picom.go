package picom

import (
	"fmt"
	"strings"
	"time"

	u "github.com/FT-Labs/physettings/utils"
	"github.com/FT-Labs/tview"
	"github.com/gdamore/tcell/v2"
)

var app *tview.Application
var pages *tview.Pages
var confirm *tview.Modal
var scriptInfo *tview.TextView
var scriptInfoLast string
var scriptRunning bool = false

var lastFocus tview.Primitive
var o1, o2 *tview.Form

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

func dropSelOpenWindowAnim(selection string, i int) {
    err := changePicomAttribute(_animation_for_open_window, selection)
    if err != nil {
        confirm.SetText(err.Error()).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorRed)
    } else {
        confirm.SetText("Open window animation changed to: " + selection).
                SetBackgroundColor(tcell.Color59).
                SetTextColor(tcell.ColorLightGreen)
    }
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
    pages.ShowPage("confirm")
}


func makeDropdown(opt string) *tview.DropDown {
    d := tview.NewDropDown()
    d.List.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune){
        scriptInfoLast = printScriptInfo(animInfo[mainText])
    })
    d.List.SetFocusFunc(func(){
        _, s := d.GetCurrentOption()
        scriptInfoLast = printScriptInfo(animInfo[s])
    })
    switch opt {
    case _animation_for_open_window:
        d.SetLabel("OPEN WINDOW ANIM : ").
            SetOptions(animOpenOpts, dropSelRofiColor).
            SetCurrentOption(0)
        d.SetFocusFunc(func(){
            scriptInfoLast = printScriptInfo("Choose window opening animation.")
        })
    case _animation_for_unmap_window:
        d.SetLabel("CLOSE WINDOW ANIM :").
            SetOptions(animCloseOpts, dropSelPowerMenuStyle).
            SetCurrentOption(0)
        d.SetFocusFunc(func() {
            scriptInfoLast = printScriptInfo("Choose window closing or unmapping animation.")
        })
    case _animation_for_next_tag:
        d.SetLabel("ANIM FOR NEXT TAG :").
            SetOptions(animNextOpts, dropSelPowerMenuType).
            SetCurrentOption(0)
        d.SetFocusFunc(func() {
            scriptInfoLast = printScriptInfo("Choose animation for incoming tag windows.")
        })
    case _animation_for_prev_tag:
    d.SetLabel("ANIM FOR PREV TAG :").
        SetOptions(animPrevOpts, dropSelPowerMenuType).
        SetCurrentOption(0)
    d.SetFocusFunc(func() {
        scriptInfoLast = printScriptInfo("Choose animation for windows that are going out from current tag.")
    })
    }
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

func printScriptInfo(s string) string {
    if !scriptRunning && s == scriptInfoLast {
        return s
    }

    run := func() {
        scriptInfo.Clear()
        arr := strings.Split(s, " ")
        for _, word := range arr {
            for i := 0; i < 20; i++ {
                time.Sleep(time.Millisecond)
                if !scriptRunning {
                    return
                }
            }
            fmt.Fprintf(scriptInfo, "%s ", word)
        }
    }

    if scriptRunning {
        scriptRunning = false
        time.Sleep(time.Millisecond * 2)
        printScriptInfo(s)
    } else {
        scriptRunning = true
        go run()
    }
    return s
}

func makeOptionsForm() *tview.Form {

    c := tview.NewCheckbox().
            SetLabel("EXPERIMENTAL BACKENDS:").
            SetChecked(u.Attrs[u.PICOM_EXPERIMENTAL] == "true").
            SetChangedFunc(checkSelExperimental)

    c.SetFocusFunc(func(){
        scriptInfoLast = printScriptInfo("Enable experimental backends in picom.\nThis will add dual-kawase blur, which is preferred.\nNote that this uses GLX backend, which doesn't work correctly in legacy hardware.")
    })

    c1 := tview.NewCheckbox().
            SetLabel("ENABLE ANIMATIONS :").
            SetChecked(picomOpts[_animations] == "true").
            SetChangedFunc(checkSelAnimConfirm)

    c1.SetFocusFunc(func(){
        scriptInfoLast = printScriptInfo("Enable or disable animations in picom")
    })

    c2 := tview.NewCheckbox().
            SetLabel("ENABLE FADING :").
            SetChecked(picomOpts[_fading] == "true").
            SetChangedFunc(checkSelFadeConfirm)

    c2.SetFocusFunc(func(){
        scriptInfoLast = printScriptInfo("Enable fading when opening - closing windows.\nWindows will go from transparent to opaque if set.")
    })

    c3 := tview.NewCheckbox().
            SetLabel("NEXT TAG FADING :").
            SetChecked(picomOpts[_enable_fading_next_tag] == "true").
            SetChangedFunc(checkSelFadeNextTag)

    c3.SetFocusFunc(func(){
        scriptInfoLast = printScriptInfo("Enable fading for incoming tag.\nNew windows that are coming from next tag will go from transparent to opaque.")
    })

    c4 := tview.NewCheckbox().
            SetLabel("PREV TAG FADING :").
            SetChecked(picomOpts[_enable_fading_prev_tag] == "true").
            SetChangedFunc(checkSelFadePrevTag)

    c4.SetFocusFunc(func(){
        scriptInfoLast = printScriptInfo("Enable fading for next tag.\nOld windows that are going out from current tag will go from opaque to transparent.")
    })

    i1 := tview.NewInputField().
            SetLabel("ANIM SPEED IN TAG :").
            SetAcceptanceFunc(tview.InputFieldFloatMaxLength(3)).
            SetPlaceholder(picomOpts[_animation_stiffness_in_tag]).
            SetFieldWidth(len(picomOpts[_animation_stiffness_in_tag]))
    i1.SetChangedFunc(func(text string){
        if len(text) > len(picomOpts[_animation_stiffness_in_tag]) {
            i1.SetFieldWidth(len(text))
        } else {
            i1.SetFieldWidth(len(picomOpts[_animation_stiffness_in_tag]))
        }
    })

    i1.SetFocusFunc(func(){
        scriptInfoLast = printScriptInfo("Set animation speed for moving or resizing windows in current tag.\nDefault value is [::b]125[::-]. Enter an integer or float number.")
    })

    i1.SetDoneFunc(func(key tcell.Key){
        switch key {
        case tcell.KeyEnter:
            err := changePicomAttribute(_animation_stiffness_in_tag, i1.GetText())

            if err != nil {
                confirm.SetText(err.Error()).
                        SetBackgroundColor(tcell.Color59).
                        SetTextColor(tcell.ColorRed)
            } else {
                confirm.SetText("Animation speed in current tag changed to: [::b]" + i1.GetText()).
                        SetBackgroundColor(tcell.Color59).
                        SetTextColor(tcell.ColorLightGreen)
            }
            pages.ShowPage("confirm")
        }
    })

    i2 := tview.NewInputField().
            SetLabel("ANIM SPEED ON TAG CHANGE :").
            SetAcceptanceFunc(tview.InputFieldFloatMaxLength(3)).
            SetPlaceholder(picomOpts[_animation_stiffness_tag_change]).
            SetFieldWidth(len(picomOpts[_animation_stiffness_tag_change]))

    i2.SetChangedFunc(func(text string){
        if len(text) > len(picomOpts[_animation_stiffness_tag_change]) {
            i2.SetFieldWidth(len(text))
        } else {
            i2.SetFieldWidth(len(picomOpts[_animation_stiffness_tag_change]))
        }
    })

    i2.SetFocusFunc(func(){
        scriptInfoLast = printScriptInfo("Set animation speed for windows transitioning between tags.\nDefault value is [::b]90[::-]. Enter an integer or float number.")
    })

    i2.SetDoneFunc(func(key tcell.Key){
        switch key {
        case tcell.KeyEnter:
            err := changePicomAttribute(_animation_stiffness_tag_change, i2.GetText())

            if err != nil {
                confirm.SetText(err.Error()).
                        SetBackgroundColor(tcell.Color59).
                        SetTextColor(tcell.ColorRed)
            } else {
                confirm.SetText("Animation speed between tags changed to: [::b]" + i2.GetText()).
                        SetBackgroundColor(tcell.Color59).
                        SetTextColor(tcell.ColorLightGreen)
            }
            pages.ShowPage("confirm")
        }
    })

    return tview.NewForm().
                SetFieldBackgroundColor(tcell.Color238).
                SetFieldTextColor(tcell.Color248).
                SetLabelColor(tcell.Color33).
                SetItemPadding(1).
                AddCheckbox(c).
                AddCheckbox(c1).
                AddCheckbox(c2).
                AddCheckbox(c3).
                AddCheckbox(c4).
                AddInputFieldItem(i1).
                AddInputFieldItem(i2)
}

func makeAnimationForm() *tview.Form {
    return tview.NewForm().
               SetItemPadding(2).
               SetLabelColor(tcell.Color111).
               SetFieldBackgroundColor(tcell.Color238).
               SetFieldTextColor(tcell.Color255).
               AddDropDown(makeDropdown(_animation_for_open_window)).
               AddDropDown(makeDropdown(_animation_for_unmap_window)).
               AddDropDown(makeDropdown(_animation_for_prev_tag)).
               AddDropDown(makeDropdown(_animation_for_next_tag))
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
            o2 = makeAnimationForm()

            o1.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
                if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
                    app.SetFocus(o2)
                    lastFocus = o2
                    return nil
                }
                return event
            })

            o2.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
                if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
                    app.SetFocus(o1)
                    lastFocus = o1
                    return nil
                }
                return event
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
		AddItem(newPrimitive("Use Tab or Up-Down keys to navigate, Shift+Tab or Left-Right to navigate between columns"), 0, 1, false).
		AddItem(newPrimitive("Type to search and Enter to select, Esc to cancel selection"), 0, 1, false)

	pages.AddPage("flex", flex, true, true).
		AddPage("confirm", confirm, true, false)

    return " 𧻓 PICOM - ANIMATIONS ", pages
}
