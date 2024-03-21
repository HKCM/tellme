---
tags: []
---
# log
---
## Example

```go
import "log/slog"
func main() {
    slog.SetLogLoggerLevel(slog.LevelDebug)
    slog.Info("hello", "username", "Mike", "age", 18)
    // 2023/08/09 20:07:51 INFO hello username=Mike age=18
    slog.Warn("hello")
    log.SetFlags(log.Ldate | log.Lmicroseconds)
    slog.Error("hello")
    // 2023/08/09 20:15:36.601583 ERROR hello
    {   //text logger
		//The second argument enables `Debug` log level.
		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(handler))
		slog.Debug("hello", "username", "Mike", "age", 18)
		// time=2024-03-21T14:41:46.650+09:00 level=DEBUG msg=hello username=Mike age=18
	}
	{   //JSON logger
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(handler))
		slog.Debug("hello", "username", "Mike", "age", 18)
		// {"time":"2024-03-21T14:41:46.650299+09:00","level":"DEBUG","msg":"hello","username":"Mike","age":18}
	}
}
```

