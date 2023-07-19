package timer

import (
	"testing"
)

func TestNewTimer(t *testing.T) {
	tmr := NewTimer(10, 2, 3, 5, 3, false)
	if tmr == nil {
		t.Error("Expected timer to not be nil")
	}
}

func TestTickBasicWork(t *testing.T) {
	tmr := NewTimer(3, 2, 3, 5, 3, false)
	state := tmr.TimerState()
	if state != STOPPED {
		t.Error("Expected state to be STOPPED. Got:", state)
	}
	state = tmr.Start()
	if state != WORK {
		t.Error("Expected state to be WORK. Got:", state)
	}
	tmr.Tick()
	if tmr.Counter() != 1 {
		t.Error("Expected counter to be 1")
	}
	state = tmr.Pause()
	if state != WORK_PAUSED {
		t.Error("Expected state to be WORK_PAUSED. Got:", state)
	}
	state = tmr.Start()
	if state != WORK {
		t.Error("Expected state to be WORK. Got:", state)
	}
}

func TestTransitions(t *testing.T) {
	maxWorkCounter := 10
	maxSbreakCounter := 2
	maxLbreakCounter := 3
	maxWorkIter := 50
	workChunk := 5
	tmr := NewTimer(
		maxWorkCounter,
		maxSbreakCounter,
		maxLbreakCounter,
		maxWorkIter,
		workChunk,
		false,
	)
	for i := 0; i < maxWorkIter-1; i++ {
		tmr.Start()
		workToPreBreakOrDone(tmr)
		if tmr.TimerState() != PRE_SBREAK && tmr.TimerState() != PRE_LBREAK {
			t.Error(
				"Expected state to be PRE_SBREAK or PRE_LBREAK. Got:",
				tmr.TimerState(),
				"on iteration",
				i,
			)
		}
		preBreakToPreWork(tmr, t)
	}
	tmr.Start()
	workToPreBreakOrDone(tmr)
	if tmr.TimerState() != DONE {
		t.Error("Expected state to be DONE. Got:", tmr.TimerState())
	}
	if tmr.WorkIter() != tmr.MaxWorkIter() {
		t.Error(
			"Expected work iterations to be",
			tmr.MaxWorkIter(),
			"Got:",
			tmr.WorkIter(),
		)
	}
}

func workToPreBreakOrDone(tmr *Timer) {
	for i := 0; i <= tmr.MaxWorkCounter(); i++ {
		tmr.Tick()
	}
}

func preBreakToPreWork(tmr *Timer, t *testing.T) {
	if (tmr.WorkIter() % tmr.WorkChunk()) != 0 {
		state := tmr.Tick()
		if state != PRE_SBREAK {
			t.Error("Expected state to be PRE_SBREAK. Got:", state)
		}
		state = tmr.Start()
		if state != SBREAK {
			t.Error("Expected state to be SBREAK. Got:", state)
		}
		state = tmr.Pause()
		if state != SBREAK_PAUSED {
			t.Error("Expected state to be SBREAK_PAUSED. Got:", state)
		}
		state = tmr.Start()
		if state != SBREAK {
			t.Error("Expected state to be SBREAK. Got:", state)
		}
		for i := 0; i <= tmr.MaxSbreakCounter(); i++ {
			tmr.Tick()
		}
		state = tmr.Tick()
		if state != PRE_WORK {
			t.Error("Expected state to be PRE_WORK. Got:", state)
		}
	} else {
		state := tmr.Tick()
		if state != PRE_LBREAK {
			t.Error("Expected state to be PRE_LBREAK. Got:", state)
		}
		state = tmr.Start()
		if state != LBREAK {
			t.Error("Expected state to be LBREAK. Got:", state)
		}
		state = tmr.Pause()
		if state != LBREAK_PAUSED {
			t.Error("Expected state to be LBREAK_PAUSED. Got:", state)
		}
		state = tmr.Start()
		if state != LBREAK {
			t.Error("Expected state to be LBREAK. Got:", state)
		}
		for i := 0; i <= tmr.MaxLbreakCounter(); i++ {
			tmr.Tick()
		}
		state = tmr.Tick()
		if state != PRE_WORK {
			t.Error("Expected state to be PRE_WORK. Got:", state)
		}
	}
}

func TestSkips(t *testing.T) {
	maxWorkCounter := 1
	maxSbreakCounter := 1
	maxLbreakCounter := 1
	maxWorkIter := 3
	workChunk := 2
	tmr := NewTimer(
		maxWorkCounter,
		maxSbreakCounter,
		maxLbreakCounter,
		maxWorkIter,
		workChunk,
		false,
	)
  if tmr.TimerState() != STOPPED {
    t.Error("Expected state to be STOPPED. Got:", tmr.TimerState())
  }
  tmr.Start()
  if tmr.TimerState() != WORK {
    t.Error("Expected state to be WORK. Got:", tmr.TimerState())
  }
  tmr.Skip()
  if tmr.TimerState() != PRE_SBREAK {
    t.Error("Expected state to be PRE_SBREAK. Got:", tmr.TimerState())
  }
  tmr.Skip()
  if tmr.TimerState() != PRE_WORK {
    t.Error("Expected state to be PRE_WORK. Got:", tmr.TimerState())
  }
  tmr.Skip()
  if tmr.TimerState() != PRE_LBREAK {
    t.Error("Expected state to be PRE_LBREAK. Got:", tmr.TimerState())
  }
  tmr.Skip()
  if tmr.TimerState() != PRE_WORK {
    t.Error("Expected state to be PRE_WORK. Got:", tmr.TimerState())
  }
  tmr.Skip()
  if tmr.TimerState() != DONE {
    t.Error("Expected state to be DONE. Got:", tmr.TimerState())
  }
}
