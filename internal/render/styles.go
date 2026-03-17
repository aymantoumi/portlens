package render

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	colorPrimary = lipgloss.Color("254") // main readable text
	colorMuted   = lipgloss.Color("67")  // secondary labels, metadata
	colorDim     = lipgloss.Color("236") // dividers, tree connectors
	colorGreen   = lipgloss.Color("77")  // ok, free ports, frontend
	colorBlue    = lipgloss.Color("75")  // info, api, compose tag
	colorAmber   = lipgloss.Color("178") // warning, suggestion, systemd tag
	colorRed     = lipgloss.Color("167") // conflict, error, unknown tag
	colorPurple  = lipgloss.Color("140") // database, cache, search
	colorOrange  = lipgloss.Color("173") // storage, infra, homelab
)

const (
	colPort  = 6
	colSvc   = 16
	colCat   = 10
	colSrc   = 8
	colBadge = 10
	colInfo  = 20
)

const (
	divider = "──────────────────────────────────────────────────────────────"
)

var (
	stylePrimary = lipgloss.NewStyle().Foreground(colorPrimary)
	styleMuted   = lipgloss.NewStyle().Foreground(colorMuted)
	styleDim     = lipgloss.NewStyle().Foreground(colorDim)
	styleGreen   = lipgloss.NewStyle().Foreground(colorGreen)
	styleBlue    = lipgloss.NewStyle().Foreground(colorBlue)
	styleAmber   = lipgloss.NewStyle().Foreground(colorAmber)
	styleYellow  = lipgloss.NewStyle().Foreground(colorAmber) // alias for amber
	styleRed     = lipgloss.NewStyle().Foreground(colorRed)
	stylePurple  = lipgloss.NewStyle().Foreground(colorPurple)
	styleOrange  = lipgloss.NewStyle().Foreground(colorOrange)
	styleWhite   = lipgloss.NewStyle().Foreground(colorPrimary) // alias for primary

	styleDivider = lipgloss.NewStyle().Foreground(colorDim)

	styleConflictRow = lipgloss.NewStyle().Background(lipgloss.Color("52"))
	styleInfoBlock   = lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(colorBlue).
				PaddingLeft(2)

	styleConflictBlock = lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(colorRed).
				PaddingLeft(2).
				Background(lipgloss.Color("52"))
)

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < 60 {
		return 80
	}
	return width
}

type TableColumns struct {
	Port  int
	Svc   int
	Cat   int
	Src   int
	Badge int
	Info  int
	Total int
}

func calcColumns() TableColumns {
	return TableColumns{
		Port:  colPort,
		Svc:   colSvc,
		Cat:   colCat,
		Src:   colSrc,
		Badge: colBadge,
		Info:  colInfo,
		Total: 0,
	}
}
