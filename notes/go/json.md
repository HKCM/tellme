---
tags: []
---
# json 
---
```go
type DbJson struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Db       string `json:"db"`
}

func getJsonConfig(filePath string) DbJson {
	// 1. 读取json文件内容
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("err1", err)
		return DbJson{}
	}
	db := new(DbJson)
	// 2. 将读取到的json文件内容，进行反序列化；将得到一个[]byte类型的切片
	err = json.Unmarshal(file, db)
	if err != nil {
		fmt.Println("err2", err)
		return DbJson{}
	}
	return *db
}

func getJsonConfigMap(filePath string) map[string]string {
	// 1. 读取json文件内容
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("err1", err)
		return nil
	}

	// 2.1 将读取到的json文件内容，进行反序列化，复制给map[string][]byte(和2中的效果是一样的)
	allConfig := make(map[string]string, 0)
	err = json.Unmarshal(file, &allConfig)
	if err != nil {
		fmt.Println("err3", err)
		return nil
	}

	return allConfig

	// 3. 循环map内容
	// for k, v := range allConfig {
	// 	fmt.Println(k, string(v)) // 值为[]byte类型，将其转为string
	// }
}
```
