---
tags: ["mysql_index"]
---
# index
---
## Example


```sql
-- 显示索引
SHOW INDEX FROM table_name\G

-- 创建索引
CREATE INDEX idx_name ON table_name (username,birth);

-- 删除索引
DROP INDEX idx_name ON table_name;

```

---
## 其他

```sql
-- 建表时指定索引
CREATE TABLE students (
  id INT PRIMARY KEY,
  name VARCHAR(50),
  age INT,
  INDEX idx_age (name,age) -- 建表时指定索引
);

-- 通过ALTER创建索引
ALTER TABLE employees ADD INDEX idx_name_age (name,age);

-- 通过ALTER删除索引
ALTER TABLE table_name DROP INDEX index_name_age;
```
