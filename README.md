# log
golang logger

## usage

```golang
import (
  "github.com/yeqown/log"
)

const (
  logPath = "./testdata"
  filename = "app"
)

func main() {
  l := NewLogger()
  // set file output, if not set, logger will only output to stderr
  l.SetFileOutput(logPath, filename)
  
  l.Info("info")
  l.Infof("%d is not equal to %d", 1, 2)
  l.Error("info")
}
```

## using preview

![screenshot](https://raw.githubusercontent.com/yeqown/log/master/screenshot.png)
