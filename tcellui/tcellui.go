package tcellui

import (
	"fmt"
	"os"
	"sync"

	"github.com/gdamore/tcell/v2"
)

type Trigger struct {
	State int
	Char  rune
}

type EventResponses map[Trigger]func()

type TcellUI struct {
	Text           string
	AppState       int
	style          tcell.Style
	screen         tcell.Screen
	eventResponses EventResponses
	sizeX          int
	sizeY          int
	prevText       string
}

func NewTcellUI(appState int) *TcellUI {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create screen: %v", err)
		os.Exit(1)
	}
	err = screen.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize screen: %v", err)
		os.Exit(1)
	}
	screen.SetStyle(
		tcell.StyleDefault.Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlack),
	)
	sizeX, sizeY := screen.Size()
	eventResponses := EventResponses{}
	tcellui := &TcellUI{
		Text:           "",
		AppState:       appState,
		style:          tcell.StyleDefault,
		screen:         screen,
		eventResponses: eventResponses,
		sizeX:          sizeX,
		sizeY:          sizeY,
		prevText:       "",
	}
	return tcellui
}

func (self *TcellUI) AddEventResponse(trigger Trigger, response func()) {
	self.eventResponses[trigger] = response
}

func (self *TcellUI) RemoveEventResponse(trigger Trigger) {
	delete(self.eventResponses, trigger)
}

// Run this in a goroutine to listen for events
func (self *TcellUI) Listen(wg *sync.WaitGroup) {
	done := make(chan struct{})
	self.sizeX, self.sizeY = self.screen.Size()
	quit := func() {
		maybePanic := recover()
		self.screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
	// self.drawText(0, 0, "Press 'q' to quit")
	go func() {
		for {
			ev := self.screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					self.screen.Fini()
					close(done)
					return
				}
				// Process other key events
				switch ev.Rune() {
				case 'q':
					self.screen.Fini()
					close(done)
					return
				default:
					r := ev.Rune()
					for k, v := range self.eventResponses {
						if k.Char == r && k.State == self.AppState {
							v()
						}
					}
				}
			case *tcell.EventResize:
				self.screen.Sync()
				self.sizeX, self.sizeY = self.screen.Size()
			}
		}
	}()
	for {
		select {
		case <-done:
			self.screen.Fini()
			wg.Done()
			return
		}
	}
}

func (self *TcellUI) drawText(x, y int, text string) {
	row := y
	col := x
	x2 := self.sizeX
	y2 := self.sizeY
	for _, r := range []rune(text) {
		self.screen.SetContent(col, row, r, nil, self.style)
		col++
		if col >= x2 || r == '\n' {
			row++
			col = x
		}
		if row > y2 {
			break
		}
	}
}

// Update screen based on what's in TcellUI.Text
func (self *TcellUI) Update() {
	if self.Text == self.prevText {
		return
	}
	self.screen.Clear()
	self.drawText(0, 0, self.Text)
	self.screen.Show()
	self.prevText = self.Text
}
