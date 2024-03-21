---
tags: []
---
# mysql
---
## Example

```bash
# 创建mysql自己的网络 在同一网络里不需要指定-p
docker network create -d bridge mysql-network
# docker network rm mysql-network
docker run --network mysql-network --name n-mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:8.0.32
# docker rm mysql
docker run -it --network mysql-network --rm mysql:8.0.32 mysql -hn-mysql -uroot -p

# 在host网络下需要指定-p
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:8.0.32
# 通过容器连接 需要使用宿主机的IP,因为在容器内localhost和127.0.0.1都不是指宿主机
docker run -it --rm mysql:8.0.32 mysql -h10.20.21.5 -uroot -p
# 通过宿主机直连需要使用127.0.0.1 不要使用localhost 否则报错
mysql -h127.0.0.1 -uroot -p 
```
