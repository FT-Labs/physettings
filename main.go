package main

import (
	"fmt"
	"strconv"

	. "github.com/FT-Labs/physettings/info"
	. "github.com/FT-Labs/physettings/keys"
	. "github.com/FT-Labs/physettings/options"
	"github.com/FT-Labs/physettings/utils"
	"github.com/FT-Labs/tview"
	"github.com/gdamore/tcell/v2"
)


type Slide func(app *tview.Application, nextSlide func()) (title string, content tview.Primitive)


var app = tview.NewApplication()

func main() {
    // Get global attributes
    utils.FetchAttributes()
    slides := []Slide{
        Cover,
        Keys,
        Options,
    }

    pages := tview.NewPages()

    info := tview.NewTextView().
        SetDynamicColors(true).
        SetRegions(true).
        SetWrap(false).
        SetHighlightedFunc(func(added, removed, remaining []string){
            pages.SwitchToPage(added[0])
        })

    // Create the pages for all slides.
	previousSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide - 1 + len(slides)) % len(slides)
		info.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}
	nextSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide + 1) % len(slides)
		info.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}
    chooseSlide := func(s int) {
        info.Highlight(strconv.Itoa(s)).
            ScrollToHighlight()
    }

    for index, slide := range slides {
        title, primitive := slide(app, nextSlide)
        pages.AddPage(strconv.Itoa(index), primitive, true, index == 0)
        fmt.Fprintf(info, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, title)
    }
    info.Highlight("0")

    layout := tview.NewFlex().
    SetDirection(tview.FlexRow).
    AddItem(pages, 0, 1, true).
    AddItem(info, 1, 1, false)

	// Shortcuts to navigate the slides.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			nextSlide()
			return nil
		} else if event.Key() == tcell.KeyCtrlP {
			previousSlide()
			return nil
		} else if event.Key() >= tcell.KeyF1 && event.Key() <= tcell.KeyF9  {
            i := event.Key() - tcell.KeyF1
            chooseSlide(int(i))
            return nil
        }
		return event
	})

	// Start the application.
	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
