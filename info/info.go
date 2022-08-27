package info

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/FT-Labs/tview"
)

const logo = `
#######################+=###-
 .=+++++++++++++++=*###=+###.
  +#############==###-*##+
    =###+=+++---:+##*-###=
     :###=+###..*##+=###-
      .*##+=###-+#=+##*.
        +##*-###+:*##+
         =###-*#####=
          :###==###:
           .*##+-*.
             +##+
              -=
`

const ftl = `
    ______________
   / ____/_  __/ /
  / /_    / / / /
 / __/   / / / /___
/_/     /_/ /_____/
`

const (
	subtitle   = `phyOS - Settings & Usage Guide`
    navigation = `[F1 .. F9]: Choose page    Ctrl+N: Next page    Ctrl+P: Previous page    Ctrl+C: Exit`
    zoom       = `Alt+Shift+J to zoom out    Alt+Shift+K to zoom in    Alt+; to cycle fonts`
	mouse      = `Or use mouse`
)


// Cover returns the cover page.
func Cover(app *tview.Application, nextSlide func()) (title string, content tview.Primitive) {
	// What's the size of the logo?
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)
	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}
	lines = strings.Split(ftl, "\n")
	ftlWidth := 0
	ftlHeight := len(lines)
	for _, line := range lines {
		if len(line) > ftlWidth {
			ftlWidth = len(line)
		}
	}

	logoBox := tview.NewTextView().
		SetTextColor(tcell.Color111).
		SetDoneFunc(func(key tcell.Key) {
			nextSlide()
		})
    ftlBox := tview.NewTextView().
        SetTextColor(tcell.Color116)
	fmt.Fprint(logoBox, logo)
    fmt.Fprint(ftlBox, ftl)

	// Create a frame for the subtitle and navigation infos.

	frame := tview.NewFrame(tview.NewBox()).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(subtitle, true, tview.AlignCenter, tcell.ColorWhite).
		AddText("", true, tview.AlignCenter, tcell.ColorWhite).
		AddText(navigation, true, tview.AlignCenter, tcell.ColorDarkMagenta).
		AddText(mouse, true, tview.AlignCenter, tcell.ColorDarkMagenta).
		AddText("", true, tview.AlignCenter, tcell.ColorDarkMagenta).
		AddText(zoom, true, tview.AlignCenter, tcell.ColorDarkCyan)

	// Create a Flex layout that centers the logo and subtitle.
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 4, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logoBox, logoWidth, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true).
		AddItem(frame, 0, 5, false).
            AddItem(tview.NewFlex().
                AddItem(tview.NewBox(), 0, 1, false).
                AddItem(ftlBox, ftlWidth, 1, true).
                AddItem(tview.NewBox(), 0, 1, false), ftlHeight, 1, false)

	return "  INFO ", flex
}
