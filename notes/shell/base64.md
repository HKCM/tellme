---
tags: ["decode"]
---
# base64
---
## Example

```shell
# base64加密
base64 filename > outputfile
cat file | base64 > outputfile

# base64解密
base64 -d file > outputfile
cat base64_file | base64 -d > outputfile
```

```bash
echo hello |base64
aGVsbG8K

echo aGVsbG8K|base64 -d
hello
```

