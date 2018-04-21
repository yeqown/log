# log
golang logger

## usage

sample-1.go
```golang
// sample-1.go
// can use `log.Info(s)` or other output function directly

import (
  "github.com/yeqown/log"
)

func main() {  
  log.Info("info")
  log.Infof("%d is not equal to %d", 1, 2)
  log.Error("info")
}
```

sample-2.go

```golang
// sample-2.go
// do `log.NewLogger` to new custom logger
// do `logger.SetFileOutput(path, filename)` to add file output

import (
  "github.com/yeqown/log"
)

const (
  logPath = "./testdata"
  filename = "app"
)

func main() {
  l := log.NewLogger()
  // set file output, if not set, logger will only output to stderr
  l.SetFileOutput(logPath, filename)
  
  l.Info("info")
  l.Infof("%d is not equal to %d", 1, 2)
  l.Error("info")
}
```

## using preview

> Note: for testing function easily, I set ticker duration as 1*time.Second* like you see in the screenshot, default is 1 *time.Minute*

![screenshot](https://raw.githubusercontent.com/yeqown/log/master/screenshot.png)
