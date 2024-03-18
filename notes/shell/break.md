---
tags: []
---
# break
---
## Example

`break`命令立即终止循环,程序继续执行循环块之后的语句,即不再执行剩下的循环

```shell
#!/usr/bin/env bash

for number in 1 2 3 
do
    for s in a b c
    do
    echo "char is $s"
    if [ "$s" = "a" ]; then
        break
    fi
    done
done
```

上面例子会打印3行a。一旦变量`$s`等于a,就会跳出内层循环,继续执行外层
