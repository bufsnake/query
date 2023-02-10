## 简介

> 逻辑查询解析器

## TODO

- [x] 为keyword添加格式检测，如ip="1.1.1.1"(正确), ip="1"(错误)，具体参考example/main.go

## 实现

> REF: https://segmentfault.com/a/1190000010998941
> 
> 基于词法分析、语法分析解析校验用户输入
> 
> 词法分析: 获取Token流
> 
> 语法分析: 解析Token流

## 关键字

> 支持的连接符

```bash
and && # 与
or ||  # 或
=      # 包含字符串
==     # 必须全部等于字符串
!=     # 不包含字符串
~=     # 正则匹配字符串
!~=    # 不包含正则匹配字符串
*=     # 通配符匹配（bleve）
!*=    # 非（通配符匹配（bleve））
```

> 关键字段

```bash
# 参考cmd/example/main.go
```

## Example TEST

### GORM

```bash
 INPUT: host="baidu.com"
FORMAT: host="baidu.com"
   SQL: `host` LIKE ?
PARAMS: [%baidu.com%]
```

### bleve

```bash
FORMAT: host*="*.baidu.com"
```

```json
{
  "query": {
    "wildcard": "moc.udiab.*",
    "field": "host"
  },
  "size": 10,
  "from": 0,
  "highlight": null,
  "fields": null,
  "facets": null,
  "explain": false,
  "sort": [
    "-_score"
  ],
  "includeLocations": false,
  "search_after": null,
  "search_before": null
}
```