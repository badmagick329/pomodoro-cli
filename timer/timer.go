package timer

import "fmt"

type TimerState int

const (
	// Timer states
	STOPPED TimerState = iota
	PRE_WORK
	WORK
	WORK_PAUSED
	PRE_SBREAK
	SBREAK
	SBREAK_PAUSED
	PRE_LBREAK
	LBREAK
	LBREAK_PAUSED
	DONE
)

func (t TimerState) String() string {
	switch t {
	case STOPPED:
		return "STOPPED"
	case PRE_WORK:
		return "PRE_WORK"
	case WORK:
		return "WORK"
	case WORK_PAUSED:
		return "WORK_PAUSED"
	case PRE_SBREAK:
		return "PRE_SBREAK"
	case SBREAK:
		return "SBREAK"
	case SBREAK_PAUSED:
		return "SBREAK_PAUSED"
	case PRE_LBREAK:
		return "PRE_LBREAK"
	case LBREAK:
		return "LBREAK"
	case LBREAK_PAUSED:
		return "LBREAK_PAUSED"
	case DONE:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}

type Timer struct {
	counter          int // in seconds
	maxWorkCounter   int // in seconds
	maxSbreakCounter int // in seconds
	maxLbreakCounter int // in seconds
	workIter         int
	maxWorkIter      int
	workChunk        int
	remainingBreaks  int
	breaksLength     int // in seconds
	AutoAdvance      bool
	timerState       TimerState
}

func NewTimer(
	maxWorkCounter,
	maxSbreakCounter,
	maxLbreakCounter,
	maxWorkIter,
	workChunk int,
	autoAdvance bool,
) *Timer {
	return &Timer{
		counter:          0,
		maxWorkCounter:   maxWorkCounter,
		maxSbreakCounter: maxSbreakCounter,
		maxLbreakCounter: maxLbreakCounter,
		workIter:         0,
		maxWorkIter:      maxWorkIter,
		workChunk:        workChunk,
		remainingBreaks:  maxWorkIter - 1,
		breaksLength:     0,
		AutoAdvance:      autoAdvance,
		timerState:       STOPPED,
	}
}

// This is responsible for updating transitions that don't require user input.
// This should be called once a second in a goroutine until the timer state is
// DONE at which point it will stop doing anything.
func (self *Timer) Tick() TimerState {
	switch self.timerState {
	case WORK:
		if self.counter >= self.maxWorkCounter {
			if self.workIter < self.maxWorkIter-1 {
				if self.workIter == 0 ||
					(self.workIter > 0 && (self.workIter+1)%self.workChunk != 0) {
					self.counter = 0
					self.workIter++
					if self.AutoAdvance {
						self.timerState = SBREAK
					} else {
						self.timerState = PRE_SBREAK
					}
				} else if (self.workIter+1)%self.workChunk == 0 {
					self.counter = 0
					self.workIter++
					if self.AutoAdvance {
						self.timerState = LBREAK
					} else {
						self.timerState = PRE_LBREAK
					}
				}
			} else {
				self.timerState = DONE
				self.workIter++
			}
		} else {
			self.counter++
		}
	case SBREAK:
		if self.counter >= self.maxSbreakCounter {
			self.counter = 0
			self.remainingBreaks--
			self.breaksLength += self.maxSbreakCounter
			if self.AutoAdvance {
				self.timerState = WORK
			} else {
				self.timerState = PRE_WORK
			}
		} else {
			self.counter++
		}
	case LBREAK:
		if self.counter >= self.maxLbreakCounter {
			self.counter = 0
			self.remainingBreaks--
			self.breaksLength += self.maxLbreakCounter
			if self.AutoAdvance {
				self.timerState = WORK
			} else {
				self.timerState = PRE_WORK
			}
		} else {
			self.counter++
		}
	}
	return self.timerState
}

func (self *Timer) Start() TimerState {
	switch self.timerState {
	case STOPPED, WORK_PAUSED, PRE_WORK:
		self.timerState = WORK
	case SBREAK_PAUSED, PRE_SBREAK:
		self.timerState = SBREAK
	case LBREAK_PAUSED, PRE_LBREAK:
		self.timerState = LBREAK
	}
	return self.timerState
}

func (self *Timer) Stop() TimerState {
	self.timerState = DONE
	return self.timerState
}

func (self *Timer) Pause() TimerState {
	switch self.timerState {
	case WORK:
		self.timerState = WORK_PAUSED
	case SBREAK:
		self.timerState = SBREAK_PAUSED
	case LBREAK:
		self.timerState = LBREAK_PAUSED
	}
	return self.timerState
}

func (self *Timer) Reset() TimerState {
	self.counter = 0
	self.workIter = 0
	self.timerState = STOPPED
	return self.timerState
}

func (self *Timer) TimerState() TimerState {
	return self.timerState
}

func (self *Timer) Counter() int {
	return self.counter
}

// Convert seconds to human readable time string.
func (self *Timer) TimeString(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%02ds", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%02dm:%02ds", seconds/60, seconds%60)
	} else {
		return fmt.Sprintf("%02dh:%02dm:%02ds", seconds/3600, seconds%3600/60, seconds%60)
	}
}

func (self *Timer) WorkIter() int {
	return self.workIter
}

func (self *Timer) MaxWorkIter() int {
	return self.maxWorkIter
}

func (self *Timer) WorkChunk() int {
	return self.workChunk
}

func (self *Timer) MaxWorkCounter() int {
	return self.maxWorkCounter
}

func (self *Timer) MaxSbreakCounter() int {
	return self.maxSbreakCounter
}

func (self *Timer) MaxLbreakCounter() int {
	return self.maxLbreakCounter
}

func (self *Timer) RemainingBreaks() int {
	return self.remainingBreaks
}

func (self *Timer) BreaksLength() int {
	return self.breaksLength
}
