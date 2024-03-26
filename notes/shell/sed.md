---
tags: []
---
# sed
---
## 常用选项
- `-n`∶ 使用安静(silent)模式。只有经过sed特殊处理的那一行(或者动作)才会被列出来。
- `-e`∶ 直接在指令列模式上进行 sed 的动作编辑
- `-f`∶ 直接将 sed 的动作写在一个档案内, -f filename 则可以执行filename 内的sed 动作
- `-r`∶ sed 的动作支援的是延伸型正规表示法的语法。(预设是基础正规表示法语法)
- `-i`∶ 直接修改读取的档案内容,而不是由屏幕输出。

## 常用命令
- `a`∶ 新增, a 的后面可以接字串,而这些字串会在新的一行出现(目前的下一行)
- `c`∶ 取代, c 的后面可以接字串,这些字串可以取代 n1,n2 之间的行！
- `d`∶ 删除
- `i`∶ 插入, i 的后面可以接字串,而这些字串会在新的一行出现(目前的上一行)；
- `p`∶ 打印,亦打印即将操作的行。通常 p 会与参数 sed -n 一起运作
- `s`∶ 取代,可以直接进行取代的工作！通常s可以搭配正则
- `=`：打印行号

## 示例
```bash
sed -n 's/test/trial/p' file                 # 替换前查看结果,只显示更改的行
sed -n 's/000/666/g;w new_test' test         # 将更改写入change.txt
sed -i.bak 's/abc/def/g' file                # 数据替换前备份
sed 's/pattern/replace_string/g' file        # 全局替换匹配的内容
sed '/^$/d' file                             # 删除空行
sed '/hello/d' file                          # 删除包含hello的行
sed -e 's/brown/green/;s/dog/cat/' data.txt  # 执行多条指令
sed '/Samantha/s/bash/csh/' /etc/passwd      # pattern匹配,仅在匹配到Samantha的行中,将bash换为zsh
```
---
## 高级用法

```shell
# 匹配字符串标记（&）
echo this is an example | sed 's/\w\+/[&]/g' 
[this] [is] [an] [example]
echo "The cat slleps in his hat" | sed 's/.at/"&"/g'
The "cat" slleps in his "hat"

# 子串匹配标记（\1）
echo this is digit 7 in a number | sed 's/digit \([0-9]\)/\1/'
this is 7 in a number 

# 引用变量
$ text=hello
$ echo hello world | sed "s/$text/HELLO/"
HELLO world 
```

## 详细示例

### 替换
```shell
$ cat data1.txt 
the first line
The quick brown fox jumps over the lazy dog. dog is cute.
The quick brown fox jumps over the lazy dog. dog is cute.
The last line.

# 简单替换(替换首个匹配)
$ sed 's/dog/cat/' data1.txt
The quick brown fox jumps over the lazy cat.
The quick brown fox jumps over the lazy cat.

# 简单替换(替换第二个匹配)
sed 's/dog/cat/2' data1.txt
The quick brown fox jumps over the lazy dog. cat is cute.
The quick brown fox jumps over the lazy dog. cat is cute.

# 简单替换(替换所有匹配)
sed 's/dog/cat/g' data1.txt
The quick brown fox jumps over the lazy cat. cat is cute.
The quick brown fox jumps over the lazy cat. cat is cute.

# 多个替换,执行多个命令用分号隔开
sed -e 's/brown/green/;s/dog/cat/;s/fox/elephant/' data1.txt
The quick green elephant jumps over the lazy cat.
The quick green elephant jumps over the lazy cat.
```

## 扩展

sed文件,如果有大量要处理的sed命令,那么将它们放进一个单独的文件中通常会更方便一些。 可以在sed命令中用-f选项来指定文件。

```shell
$ cat script1.sed 
s/brown/green/
s/fox/elephant/
s/dog/cat/

$ sed -f script1.sed data1.txt 
The quick green elephant jumps over the lazy cat.
The quick green elephant jumps over the lazy cat.
```

### 指定位置替换
```shell
sed 's/dog/cat/' data1.txt # 只作用于每行的第一次
sed '2s/dog/cat/' data1.txt # 只作用于第二行
sed '2,3s/dog/cat/' data1.txt # 作用于2-3行
sed '2,$s/dog/cat/' data1.txt # 第二行至末行
```

### 文本模式过滤

文本模式过滤模式在前方写入pattern,会匹配具有这个pattern的行

```shell
# 在匹配到Samantha的行中,将bash换为zsh
sed '/Samantha/s/bash/csh/' /etc/passwd # 在匹配到Samantha的行中,将bash换为zsh
```

在单行中执行多条命令
```shell
sed '2{s/Two/2/;s/test/real/}' data2.txt
sed '3,${s/Two/2/;s/test/real/}' data2.txt
```

### 删除行
```shell
sed '3d' data6.txt # 删除单行
sed '3,$d' data2.txt # 删除区间
sed '/number 1/d' data6.txt # 模式匹配删除
```

### 插入和追加
```shell
sed '1i\Test Line 1' data2.txt 
sed '$a\Test Last Line ' data2.txt 

# 通过读取文件的形式追加
$ cat data12.txt
This is an added line.
This is the second added line.
$ sed '3r data12.txt' data6.txt 
This is line number 1.
This is line number 2.
This is line number 3.
This is an added line.
This is the second added line. 
This is line number 4.
```

### 修改整行
```shell
# 以模式匹配方式整行替换
sed '/One/c\new line' data2.txt 
# 以行号方式
sed '1c\new line' data2.txt 
```

### 映射转换
```shell
$ cat data3.txt 
This is line number 1.
This is line number 2.
This is line number 3.
This is line number 4.
$ sed 'y/123/456/' data3.txt
This is line number 4.
This is line number 5.
This is line number 6.
This is line number 4.
```

sed编辑器不会修改原始文件。你删除的行只是从sed编辑器的输出中消失了。原始文件仍然包含那些“删掉的”行。

```shell
sed -n '$p' data2.txt # 打印尾行 tail -n 1
sed '/./,/^$/!d' data8.txt # 无论文件的数据行之间出现了多少空白行,在输出中只会在行间保留一个空白行
sed '=' data4 | sed 'N;s/\n/ /' # 给文件添加行号
```

### 删除行首空格
```shell
sed 's/^[ ]*//g' filename
sed 's/^ *//g' filename
sed 's/^[[:space:]]*//g' filename
```

### 删除HTML标签
```shell
# 排除大于号,否则会进行贪婪匹配,会删除类似<b>first</b>这样的加粗文本
sed 's/<[^>]*>//g ; /^$/d' data11.txt
```

## 问题

**sed: -i may not be used with stdin**

MAC系统上sed使用时需要 
```shell
sed -i '' 's/a/b/' filename
```
