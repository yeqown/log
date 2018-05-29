# log
golang logger based `log`

## Doc
ref to: [https://godoc.org/github.com/yeqown/log](https://godoc.org/github.com/yeqown/log)

## Usage

sample-1.go
```golang
import (
  "github.com/yeqown/log"
)

func main() {
  intptr := new(int)
  *intprt = 9999

  struct_var := struct {
    Name string
    Age  int
  }{"Tonn", 24}

  // to set file output for default logger
  SetFileOutput("/path/to/logfile", "default")

  // also support Debug, Warn, Fatal, Error
  log.Info("this is a struct var: ", struct_var)
  log.Info("this is a int ptr and var: ", a, *a)
  log.Infof("%d is not equal to %d", 1, 2)
}
```

sample-2.go

```golang
import (
  "github.com/yeqown/log"
)

func main() {
  // to make self logger
  l := log.NewLogger()

  intptr := new(int)
  *intprt = 9999

  struct_var := struct {
    Name string
    Age  int
  }{"Tonn", 24}

  // to set file output for default logger
  l.SetFileOutput("/path/to/logfile", "app")

  // also support Debug, Warn, Fatal, Error
  l.Info("this is a struct var: ", struct_var)
  l.Info("this is a int ptr and var: ", a, *a)
  l.Infof("%d is not equal to %d", 1, 2)
}
```

## Using preview

> Note: this screenshot is log_test.go output's screenshot.

![screenshot](https://raw.githubusercontent.com/yeqown/log/master/screenshot.png)
