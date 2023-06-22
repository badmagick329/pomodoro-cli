package main

import (
	"os"
	"pomodoro/runner"
	"sync"

	"github.com/akamensky/argparse"
)

const CONFIG_PATH = "go_pomodoro_config.json"

func main() {
	parser := argparse.NewParser("pomodoro", "A Simple and Customisable CLI Pomodoro timer")
	var workTime *int = parser.Int("w", "work", &argparse.Options{Required: false, Help: "Work time in minutes"})
	var breakTime *int = parser.Int("b", "break", &argparse.Options{Required: false, Help: "Break time in minutes"})
	var longBreakTime *int = parser.Int("l", "long-break", &argparse.Options{Required: false, Help: "Long break time in minutes"})
	var longBreakInterval *int = parser.Int("i", "interval", &argparse.Options{Required: false, Help: "Number of pomodoros before a long break"})
	var totalPomodoros *int = parser.Int("n", "number", &argparse.Options{Required: false, Help: "Total number of pomodoros"})
	var autoStart *bool = parser.Flag("a", "auto-start", &argparse.Options{Required: false, Help: "When a timer finishes, auto start the next timer"})
	parser.Parse(os.Args)
	cfg, _ := runner.NewConfig(CONFIG_PATH)
	cfg.ReadArgs(*workTime, *breakTime, *longBreakTime, *longBreakInterval, *totalPomodoros, *autoStart)
	// cfg := runner.TestConfig()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go runner.Run(&wg, &cfg)
	wg.Wait()
}
