---
tags: ["getkey"]
---
# redis-key
---
## Example

```bash
TYPE key # 获取key的类型
KEYS my* # 获取所有以my开头的ky
DEL key  # 删除key

GET mykey          # 获取string类型的 mykey的值
LRANGE mylist 0 -1 # 获取List类型的 mylist的值
SMEMBERS myset     # 获取Set类型 myset的值
HGETALL user       # 获取hash类型的 user的值
HGET user name     # 获取hash类型的 user的name 字段的值
ZRANGE zlist 0 -1  # 获取Zset类型的 zlist的值
```

命令查询：https://www.redis.com.cn/commands
