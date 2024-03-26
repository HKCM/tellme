---
tags: ["时间格式化","时间加减","timeformat"]
---
# time 
---
```go
func main() {
	now := time.Now() // 获取当前时间
	fmt.Printf("current time:%v\n", now)
	year := now.Year()     // 年
	month := now.Month()   // 月
	day := now.Day()       // 日
	hour := now.Hour()     // 小时
	minute := now.Minute() // 分钟
	second := now.Second() // 秒
	fmt.Println(year, month, day, hour, minute, second)
	t := 3
	targetTime := now.Add(time.Duration(t) * time.Hour) // 时间加减
	fmt.Printf("current time:%v\n", targetTime.Format("2006-01-02T15:04:05Z"))
}
```
