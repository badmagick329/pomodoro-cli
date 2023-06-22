## Pomodoro-CLI

A simple and customisable CLI pomodoro timer

https://github.com/badmagick329/pomodoro-cli/assets/63713349/ac1f3ceb-f13f-43cf-9684-ac17becacb74

### Usage

```
pomodoro [-h|--help] [-w|--work <integer>] [-b|--break <integer>]
[-l|--long-break <integer>] [-i|--interval <integer>] [-n|--number <integer>]
[-a|--auto-start]

Arguments:

  -h  --help        Print help information
  -w  --work        Work time in minutes (default: 25)
  -b  --break       Break time in minutes (default: 5)
  -l  --long-break  Long break time in minutes (default: 15)
  -i  --interval    Number of pomodoros before a long break (default: 4)
  -n  --number      Total number of pomodoros (default: 8)
  -a  --auto-start  When a timer finishes, auto start the next timer (default: false)
```

### Config

The `go_pomodoro_config.json` should be placed in the same directory as the
binary. If one is not found it will be created with default values. This can
also be used to specify paths for other mp3 sound files to play when work or
break finishes.

### Alarm

By default `bell.mp3` will be used as an alarm sound when a timer finishes. 
