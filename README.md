# autopprof

Pprof made easy at development time.

## Guide

Add autopprof.Capture to your main function.

```go
import "github.com/rakyll/autopprof"

autopprof.Capture(autopprof.CPUProfile{
    Duration: 15 * time.Second,
})
```

Run your program and send SIGQUIT to the process
(or press CTRL+\\).

Profile capturing will start. Pprof UI will be started
once capture is completed.
