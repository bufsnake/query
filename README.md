## 简介

> 逻辑查询解析器

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
```

> 关键字段

```bash
# 参考cmd/example/main.go
```

## TEST

```bash
input: protocol=="https" && "127.0.0.1" and ip="1" and (title = "1"|| title="2")

output:
   SQL: `protocol` = ? AND (`ip` LIKE ? OR `ipx` LIKE ? OR `port` LIKE ? OR `protocol` LIKE ? OR `url` LIKE ? OR `location` LIKE ? OR `title` LIKE ?) AND `ip` LIKE ? AND (`title` LIKE ? OR `title` LIKE ?)
PARAMS: [https %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %1% %1% %2%]
FORMAT: protocol=="https" && "127.0.0.1" && ip="1" && (title="1" || title="2")
```