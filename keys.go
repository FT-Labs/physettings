package main

import (
	"strings"
    "os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var appW, appH int

func SetSize(w, h int) {
    appW = w
    appH = h
}

func Keys(nextSlide func()) (title string, content tview.Primitive){
    out, err := exec.Command("manfilter", "dwm").Output()

    if err != nil {
        panic("Can't read data from dwm manual page")
    }

	table := tview.NewTable().
		SetBorders(true)

    keys := strings.Split(string(out), ";")
	cols, rows := 2, len(keys)/2+1
    table.SetFixed(cols, rows)
    table.SetEvaluateAllRows(true)
    table.SetBorderAttributes(tcell.AttrBold)
	word := 0
    maxLen := 0
    for i := 0; i < len(keys) - 1; i+=2 {
        if len(keys[i]) + len(keys[i+1]) > maxLen {
            maxLen = len(keys[i]) + len(keys[i+1])
        }
    }
    color := tcell.ColorBlue
    table.SetCell(0, 0,
        tview.NewTableCell("[::b]KEY BINDING").
            SetTextColor(color).
            SetAlign(tview.AlignCenter))

    table.SetCell(0, 1,
        tview.NewTableCell("[::b]ACTION").
            SetTextColor(color).
            SetAlign(tview.AlignCenter))

	for r := 1; r < rows; r++ {
		for c := 0; c < cols; c++ {
            if c == 0 {
                color = tcell.ColorWhite
            } else {
                color = tcell.ColorLightGreen
            }
			table.SetCell(r, c,
				tview.NewTableCell(keys[word]).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			word = (word + 1) % len(keys)
		}
	}
    _, _, w, h := table.GetRect()
    table.SetRect(appW - w/2, 0, w, h)
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})

    flex := tview.NewFlex().
        SetDirection(tview.FlexColumn).
        AddItem(tview.NewBox(), 0, 1, false).
        AddItem(tview.NewFlex().
            SetDirection(tview.FlexRow).
            AddItem(table, 0, 1, true), 0, 6, true)

    return "KEY SHEET", flex
}
