---
tags: ["create_table","create_database","create_index"]
---
# <Title>
---
## Example

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS mydatabase DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- 创建表
CREATE TABLE IF NOT EXISTS users (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT 'user id',
    `username` VARCHAR(50) NOT NULL COMMENT 'user name',
    `birthdate` DATE,
    `type` tinyint NOT NULL DEFAULT '0',
    `is_active` BOOLEAN DEFAULT TRUE
) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- 创建索引
CREATE INDEX idx_username ON users (username);

-- 插入数据
insert into users (username,birthdate) VALUES 
("Ali","1994-01-01"),
("Bob","1995-01-01"),
("Cat","1996-01-01");
```
