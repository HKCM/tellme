---
tags: []
---
# redis
---
## Example

```bash
# 创建redis自己的网络
docker network create -d bridge redis-network
# docker network rm redis-network
docker run --network redis-network --name my-redis -d redis:7.0
# docker rm my-redis
docker run -it --network redis-network --rm redis:7.0 redis-cli -h my-redis
```

命令查询：https://www.redis.com.cn/commands
