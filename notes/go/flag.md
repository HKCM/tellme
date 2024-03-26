---
tags: []
---
# flag
---
```go
var (
	name string
	age int
	married bool
	delay time.Duration
)
func init(){
	flag.StringVar(&name, "name", "张三", "姓名")
	flag.IntVar(&age, "age", 18, "年龄")
	flag.BoolVar(&married, "married", false, "婚否")
	flag.DurationVar(&delay, "delay", 0, "延迟的时间间隔")
	flag.Parse() //解析命令行参数
}

// go run main.go --help
// go run main.go --name aaa --age 12 --married=true --delay 75s	
```
