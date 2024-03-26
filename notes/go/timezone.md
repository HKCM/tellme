---
tags: ["timezone","时区"]
---
# time
---
## Example

```go
layout := "2006-01-02T15:04:05Z"
// 方式一
jstZone := time.FixedZone("JST", 9*3600) // 东九 日本时间
fmt.Println(time.Now().In(jstZone).Format(layout))
	
// 方式二
tokyo, _ := time.LoadLocation("Asia/Tokyo")
fmt.Println(time.Now().In(tokyo).Format(layout))
```

