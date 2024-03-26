---
tags: []
---
# csv 
---
# 单行读取
```go
f, err := os.Open("numbers.csv")
if err != nil {
    log.Fatal(err)
}
r := csv.NewReader(f)
r.Comma = ';'
for {
    record, err := r.Read()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    for value := range record {
        fmt.Printf("%s\n", record[value])
    }
}
```

# 单行写入
```go
records := [][]string{
    {"first_name", "last_name", "occupation"},
    {"John", "Doe", "gardener"},
    {"Lucy", "Smith", "teacher"},
    {"Brian", "Bethamy", "programmer"},
}
f, err = os.Create(dbFile)
if err != nil {
    panic(err)
}
defer f.Close()
w := csv.NewWriter(f)
defer w.Flush()
for record := range records {
    err = w.Write(record})
    if err != nil {
        panic(err)
    }
}
```

# 一次性读取所有
```go
func main() {
    records, err := readData("users.csv")
    if err != nil {
        log.Fatal(err)
    }
    for _, record := range records {
        user := User{
            firstName:  record[0],
            lastName:   record[1],
            occupation: record[2],
        }
        fmt.Printf("%s %s is a %s\n", user.firstName, user.lastName,user.occupation)
    }
}
func readData(fileName string) ([][]string, error) {
    f, err := os.Open(fileName)
    if err != nil {
        return [][]string{}, err
    }
    defer f.Close()
    r := csv.NewReader(f)
    // skip first line
    if _, err := r.Read(); err != nil {
        return [][]string{}, err
    }
    records, err := r.ReadAll()
    if err != nil {
        return [][]string{}, err
    }
    return records, nil
}
```

# 一次性写入
```go
records := [][]string{
    {"first_name", "last_name", "occupation"},
    {"John", "Doe", "gardener"},
    {"Lucy", "Smith", "teacher"},
    {"Brian", "Bethamy", "programmer"},
}
f, err := os.Create("users.csv")
defer f.Close()
if err != nil {

    log.Fatalln("failed to open file", err)
}
w := csv.NewWriter(f)
defer w.Flush()
err = w.WriteAll(records) // calls Flush internally
if err != nil {
    log.Fatal(err)
}
```
