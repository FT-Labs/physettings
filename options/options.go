package options

import (
	"fmt"

	utils "github.com/FT-Labs/physettings/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application
var pages *tview.Pages
var confirm *tview.Modal
var scriptInfo *tview.TextView

var lastFocus tview.Primitive
var lastFocusIndex int
var o1, o2 *tview.Form

func dropSelRofiColor(selection string, i int) {
    utils.SetRofiColor(selection)
    confirm.SetText("Rofi colorscheme changed to: " + selection).
            SetBackgroundColor(tcell.Color59)
    lastFocusIndex = i
    pages.ShowPage("confirm")
}

func dropSelPowerMenuType(selection string, i int) {
    err := utils.SetAttribute(utils.POWERMENU_TYPE, selection)
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
    err := utils.SetAttribute(utils.POWERMENU_STYLE, selection)
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
    err := utils.RunScript(utils.POS_GRUB_CHOOSE_THEME)
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
    err := utils.RunScript(utils.POS_SDDM_CHOOSE_THEME)
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
    err := utils.RunScript(utils.POS_MAKE_BAR)
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
    if opt == utils.ROFI_COLOR {
        return utils.ROFI_COLOR + " : ",
                    utils.RofiColors,
                    0,
                    dropSelRofiColor
    } else if opt == utils.POWERMENU_STYLE {
        return       utils.POWERMENU_STYLE + " : ",
                     utils.PowerMenuStyles,
                     0,
                     dropSelPowerMenuStyle
    }// else if opt == POWERMENU_TYPE {
    return utils.POWERMENU_TYPE + " : ",
           utils.PowerMenuTypes,
           0,
           dropSelPowerMenuType
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

type Button struct {
    *tview.Button
}

func (b Button) GetFieldWidth() int {
    return len(b.GetLabel())
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
                AddDropDown(makeDropdown(utils.ROFI_COLOR)).
                AddDropDown(makeDropdown(utils.POWERMENU_TYPE)).
                AddDropDown(makeDropdown(utils.POWERMENU_STYLE))
}

func makeScriptsForm() *tview.Form {
    b := Button{tview.NewButton(utils.POS_GRUB_CHOOSE_THEME).SetSelectedFunc(buttonSelGrubTheme)}
    b.SetFocusFunc(func(){
        fmt.Fprintf(scriptInfo, "%s", utils.ScriptInfo[utils.POS_GRUB_CHOOSE_THEME])
    })
    return tview.NewForm().
               SetItemPadding(3).
               SetFieldBackgroundColor(tcell.Color16).
               SetFieldTextColor(tcell.Color231).
               AddFormItem(b).
               AddFormItem(Button{tview.NewButton(utils.POS_SDDM_CHOOSE_THEME).SetSelectedFunc(buttonSelSddmTheme)}).
               AddFormItem(Button{tview.NewButton("MAKE STATUSBAR").SetSelectedFunc(buttonSelMakeBar)})
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
		AddItem(newPrimitive("Press Tab to navigate in current column, Shift+Tab to switch between columns"), 0, 1, false).
		AddItem(newPrimitive("Enter to select (type to search, or use arrow keys)"), 0, 1, false)

	pages.AddPage("flex", flex, true, true).
		AddPage("confirm", confirm, true, false)

    return "OPTIONS", pages
}
