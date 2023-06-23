package runner

import (
	"pomodoro/player"
	"pomodoro/tcellui"
	"pomodoro/timer"
	"sync"
	"time"
)

type Markers struct {
	WorkChar  string
	BreakChar string
	EmptyChar string
}

func Run(wg *sync.WaitGroup, cfg *Config) {
	p := player.NewPlayer(cfg.WorkSoundPath, cfg.BreakSoundPath)
	tmr := timer.NewTimer(
		cfg.WorkTime,
		cfg.BreakTime,
		cfg.LongBreakTime,
		cfg.TotalPomodoros,
		cfg.LongBreakInterval,
		cfg.AutoStart,
	)
	markers := &Markers{
		WorkChar:  cfg.WorkChar,
		BreakChar: cfg.BreakChar,
		EmptyChar: cfg.EmptyChar,
	}
	ui := tcellui.NewTcellUI(0)
	updateText(ui, tmr, markers)
	addEventResponses(ui, tmr)
	go ui.Listen(wg)
	go updateLoop(tmr, ui, markers)
	go soundLoop(tmr, &p)
}

func updateLoop(tmr *timer.Timer, ui *tcellui.TcellUI, markers *Markers) {
	updateRate := 100 * time.Millisecond
	threshold := 10
	counter := 0
	for {
		ui.Update()
		if counter >= threshold {
			tmr.Tick()
			counter = 0
		}
		updateText(ui, tmr, markers)
		ui.AppState = int(tmr.TimerState())
		time.Sleep(updateRate)
		counter += 1
	}
}

func soundLoop(tmr *timer.Timer, p *player.Player) {
	updateRate := 100 * time.Millisecond
	prevState := tmr.TimerState()
	for {
		if prevState == timer.WORK {
			state := tmr.TimerState()
			autoAdvChecks := tmr.AutoAdvance && (state == timer.SBREAK ||
				state == timer.LBREAK ||
				state == timer.DONE)
			manualAdvChecks := !tmr.AutoAdvance && (state == timer.PRE_SBREAK ||
				state == timer.PRE_LBREAK ||
				state == timer.DONE)
			if autoAdvChecks || manualAdvChecks {
				prevState = tmr.TimerState()
				p.PlayWork()
			}
		} else if prevState == timer.SBREAK || prevState == timer.LBREAK {
			autoAdvChecks := tmr.AutoAdvance && tmr.TimerState() == timer.WORK
			manualAdvChecks := !tmr.AutoAdvance && tmr.TimerState() == timer.PRE_WORK
			if autoAdvChecks || manualAdvChecks {
				prevState = tmr.TimerState()
				p.PlayBreak()
			}
		} else {
			prevState = tmr.TimerState()
		}
		time.Sleep(updateRate)
	}
}

func updateText(ui *tcellui.TcellUI, tmr *timer.Timer, m *Markers) {
	ui.Text = pomoDoroString(tmr, m.WorkChar, m.EmptyChar)
	ui.Text += breakString(tmr, m.BreakChar)
	state := tmr.TimerState()
	ui.Text += timerText(state, tmr)
	ui.Text += keyText(state)
}

func keyText(state timer.TimerState) string {
	switch state {
	case timer.WORK, timer.SBREAK, timer.LBREAK:
		return "Press 'p' to pause, 'q' to quit"
	case timer.WORK_PAUSED, timer.SBREAK_PAUSED, timer.LBREAK_PAUSED:
		return "Press 's' to resume, 'q' to quit"
	case timer.PRE_SBREAK, timer.PRE_LBREAK:
		return "Press 's' to start break, 'q' to quit"
	case timer.PRE_WORK, timer.STOPPED:
		return "Press 's' to start work, 'q' to quit"
	case timer.DONE:
		return "Press 'q' to quit"
	default:
		return ""
	}
}

func timerText(state timer.TimerState, tmr *timer.Timer) string {
	switch state {
	case timer.PRE_WORK, timer.STOPPED:
		return "Work: " + tmr.TimeString(tmr.MaxWorkCounter()) + "\n\n"
	case timer.WORK, timer.WORK_PAUSED:
		return "Work: " + tmr.TimeString(
			tmr.MaxWorkCounter()-tmr.Counter(),
		) + "\n\n"
	case timer.PRE_SBREAK:
		return "Short Break: " + tmr.TimeString(tmr.MaxSbreakCounter()) + "\n\n"
	case timer.SBREAK, timer.SBREAK_PAUSED:
		return "Short Break: " + tmr.TimeString(
			tmr.MaxSbreakCounter()-tmr.Counter(),
		) + "\n\n"
	case timer.PRE_LBREAK:
		return "Long Break: " + tmr.TimeString(tmr.MaxLbreakCounter()) + "\n\n"
	case timer.LBREAK, timer.LBREAK_PAUSED:
		return "Long Break: " + tmr.TimeString(
			tmr.MaxLbreakCounter()-tmr.Counter(),
		) + "\n\n"
	case timer.DONE:
		totalWork := tmr.MaxWorkCounter() * tmr.MaxWorkIter()
		totalBreaks := tmr.BreaksLength()
		text := "Done! You worked for " + tmr.TimeString(totalWork)
		text += " and took breaks for " + tmr.TimeString(totalBreaks) + "\n\n"
		return text
	default:
		return ""
	}
}

// Respond to events and update ui text in these functions
func addEventResponses(ui *tcellui.TcellUI, tmr *timer.Timer) {
	key := tcellui.Trigger{State: int(timer.STOPPED), Char: 's'}
	startFunc := func() {
		tmr.Start()
	}
	pauseFunc := func() {
		tmr.Pause()
	}
	ui.AddEventResponse(key, startFunc)
	key = tcellui.Trigger{State: int(timer.PRE_WORK), Char: 's'}
	ui.AddEventResponse(key, startFunc)
	key = tcellui.Trigger{State: int(timer.WORK), Char: 'p'}
	ui.AddEventResponse(key, pauseFunc)
	key = tcellui.Trigger{State: int(timer.WORK_PAUSED), Char: 's'}
	ui.AddEventResponse(key, startFunc)
	key = tcellui.Trigger{State: int(timer.PRE_SBREAK), Char: 's'}
	ui.AddEventResponse(key, startFunc)
	key = tcellui.Trigger{State: int(timer.SBREAK), Char: 'p'}
	ui.AddEventResponse(key, pauseFunc)
	key = tcellui.Trigger{State: int(timer.SBREAK_PAUSED), Char: 's'}
	ui.AddEventResponse(key, startFunc)
	key = tcellui.Trigger{State: int(timer.PRE_LBREAK), Char: 's'}
	ui.AddEventResponse(key, startFunc)
	key = tcellui.Trigger{State: int(timer.LBREAK), Char: 'p'}
	ui.AddEventResponse(key, pauseFunc)
	key = tcellui.Trigger{State: int(timer.LBREAK_PAUSED), Char: 's'}
	ui.AddEventResponse(key, startFunc)
}

func pomoDoroString(tmr *timer.Timer, wc, ec string) string {
	pomodoros := ""
	for i := 0; i < tmr.WorkIter(); i++ {
		pomodoros += wc + " "
	}
	remaining := ""
	for i := 0; i < tmr.MaxWorkIter()-tmr.WorkIter(); i++ {
		pomodoros += ec + " "
	}
	return "Pomodoros: " + pomodoros + remaining + "\n"
}

func breakString(tmr *timer.Timer, bc string) string {
	remaining := ""
	for i := 0; i < tmr.RemainingBreaks(); i++ {
		remaining += bc + " "
	}
	if remaining == "" {
		return ""
	}
	return "Breaks Left: " + remaining + "\n"
}
