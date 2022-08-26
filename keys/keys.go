package keys

import (
	"os/exec"
	"strings"
	"github.com/gdamore/tcell/v2"
	"github.com/FT-Labs/tview"
)

func Keys(app *tview.Application, nextSlide func()) (title string, content tview.Primitive){
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
	word := 0
    maxC1, maxC2 := 0, 0
    for i := range keys {
        if i%2 == 0 && maxC1 < len(keys[i]) {
            maxC1 = len(keys[i])
        } else if maxC2 < len(keys[i]) {
            maxC2 = len(keys[i])
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
                    SetExpansion(2).
					SetAlign(tview.AlignCenter))
			word = (word + 1) % len(keys)
		}
	}
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
            AddItem(table, 0, 1, true), maxC1 + maxC2 + 3, 2, true).
        AddItem(tview.NewBox(), 0, 1, false)

    return "KEY SHEET", flex
}
