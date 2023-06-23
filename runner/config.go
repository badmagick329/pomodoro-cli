package runner

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const (
	DEFAULT_WORK_TIME           = 25 * 60
	DEFAULT_BREAK_TIME          = 5 * 60
	DEFAULT_LONG_BREAK_TIME     = 15 * 60
	DEFAULT_TOTAL_POMODOROS     = 8
	DEFAULT_LONG_BREAK_INTERVAL = 4
	DEFAULT_SOUND_PATH          = "bell.mp3"
)

const (
	TEST_WORK_TIME           = 2
	TEST_BREAK_TIME          = 1
	TEST_LONG_BREAK_TIME     = 3
	TEST_TOTAL_POMODOROS     = 5
	TEST_LONG_BREAK_INTERVAL = 2
)

type Config struct {
	WorkSoundPath     string `json:"work_mp3"`
	BreakSoundPath    string `json:"break_mp3"`
	WorkTime          int    `json:"work_time"`
	BreakTime         int    `json:"break_time"`
	LongBreakTime     int    `json:"long_break_time"`
	LongBreakInterval int    `json:"long_break_interval"`
	AutoStart         bool   `json:"auto_start"`
	TotalPomodoros    int    `json:"total_pomodoros"`
}

func NewConfig(configPath string) (cfg Config, errs []error) {
	cfg = Config{
		WorkSoundPath:     DEFAULT_SOUND_PATH,
		BreakSoundPath:    DEFAULT_SOUND_PATH,
		WorkTime:          DEFAULT_WORK_TIME,
		BreakTime:         DEFAULT_BREAK_TIME,
		LongBreakTime:     DEFAULT_LONG_BREAK_TIME,
		LongBreakInterval: DEFAULT_LONG_BREAK_INTERVAL,
		AutoStart:         false,
		TotalPomodoros:    DEFAULT_TOTAL_POMODOROS,
	}
	readErrs := cfg.createOrRead(configPath)
	if len(readErrs) > 0 {
		errs = append(errs, readErrs...)
	}
	return
}

func (self *Config) TestMode() {
	self.WorkTime = TEST_WORK_TIME
	self.BreakTime = TEST_BREAK_TIME
	self.LongBreakTime = TEST_LONG_BREAK_TIME
	self.LongBreakInterval = TEST_LONG_BREAK_INTERVAL
	self.TotalPomodoros = TEST_TOTAL_POMODOROS
}

func (self *Config) createOrRead(configPath string) (errs []error) {
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.Create(configPath)
			self.toMinutes()
			jsonCfg, err := json.MarshalIndent(self, "", "  ")
			if err != nil {
				self.toSeconds()
				errs = append(errs, err)
				return
			}
			err = ioutil.WriteFile(configPath, jsonCfg, 0644)
			if err != nil {
				self.toSeconds()
				errs = append(errs, err)
				return
			}
		} else {
			errs = append(errs, err)
			return
		}
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&self)
	if err == nil {
		self.toSeconds()
	} else {
		errs = append(errs, err)
	}
	validErr := self.validateSoundPaths(configPath)
	if validErr != nil {
		errs = append(errs, validErr)
	}
	return
}

func (self *Config) ReadArgs(
	workTime, breakTime, longBreakTime, longBreakInterval, totalPomodoros int,
	autoStart bool,
) {
	if workTime > 0 {
		self.WorkTime = workTime * 60
	}
	if breakTime > 0 {
		self.BreakTime = breakTime * 60
	}
	if longBreakTime > 0 {
		self.LongBreakTime = longBreakTime * 60
	}
	if longBreakInterval > 0 {
		self.LongBreakInterval = longBreakInterval
	}
	if totalPomodoros > 0 {
		self.TotalPomodoros = totalPomodoros
	}
	if autoStart {
		self.AutoStart = autoStart
	}
}

func (self *Config) validateSoundPaths(configPath string) (err error) {
	invalidWorkPath := false
	invalidBreakPath := false
	defaultExists := false
	_, e := os.Stat(DEFAULT_SOUND_PATH)
	defaultExists = !os.IsNotExist(e)
	if self.WorkSoundPath != "" {
		_, e = os.Stat(self.WorkSoundPath)
		invalidWorkPath = os.IsNotExist(e)
	}
	if self.BreakSoundPath != "" {
		_, e = os.Stat(self.BreakSoundPath)
		invalidBreakPath = os.IsNotExist(e)
	}
	if !(invalidWorkPath && invalidBreakPath) {
		return
	}
	configMsg := "Ensure the path is correct in the " + configPath + " file.\n"
	if (invalidWorkPath && invalidBreakPath) && defaultExists {
		errMsg := "Sound files " + self.WorkSoundPath + " and " +
			self.BreakSoundPath + " not found.\n"
		errMsg += configMsg
		errMsg += "Defaulting to " + DEFAULT_SOUND_PATH + "...\n"
		err = errors.New(errMsg)
		self.WorkSoundPath = DEFAULT_SOUND_PATH
		self.BreakSoundPath = DEFAULT_SOUND_PATH
		return
	}
	if (invalidWorkPath && invalidBreakPath) && !defaultExists {
		errMsg := "Sound files " + self.WorkSoundPath + " and " +
			self.BreakSoundPath + " not found.\n"
		errMsg += configMsg
		errMsg += "Default sound file " + DEFAULT_SOUND_PATH + " not found either.\n"
		err = errors.New(errMsg)
		self.WorkSoundPath = ""
		self.BreakSoundPath = ""
		return
	}
	if invalidWorkPath {
		self.WorkSoundPath = self.BreakSoundPath
		return
	} else {
		self.BreakSoundPath = self.WorkSoundPath
		return
	}
}

func (self *Config) toMinutes() {
	self.WorkTime /= 60
	self.BreakTime /= 60
	self.LongBreakTime /= 60
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (self *Config) toSeconds() {
	self.WorkTime = maxInt(self.WorkTime, 1)
	self.BreakTime = maxInt(self.BreakTime, 1)
	self.LongBreakTime = maxInt(self.LongBreakTime, 1)
	self.WorkTime *= 60
	self.BreakTime *= 60
	self.LongBreakTime *= 60
}
