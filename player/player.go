package player

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	workSoundPath  string
	breakSoundPath string
}

func NewPlayer(workSoundPath, breakSoundPath string) (p Player) {
	p.workSoundPath = workSoundPath
	p.breakSoundPath = breakSoundPath
	return
}

func playSound(soundPath string) error {
	f, err := os.Open(soundPath)
	if err != nil {
		return err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	return nil
}

func (self *Player) PlayBreak() {
	if self.breakSoundPath == "" {
		return
	}
	playSound(self.breakSoundPath)
}

func (self *Player) PlayWork() {
	if self.workSoundPath == "" {
		return
	}
	playSound(self.workSoundPath)
}
