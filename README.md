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
 INPUT: IP="127.0.0.1"||127.0.0.1 || 1234 || HOST=1
FORMAT: ip="127.0.0.1" || "127.0.0.1" || "1234" || Host="1"
   SQL: `ip` LIKE ? OR (`ip` LIKE ? OR `ipx` LIKE ? OR `port` LIKE ? OR `protocol` LIKE ? OR `url` LIKE ? OR `location` LIKE ? OR `title` LIKE ? OR `Host` LIKE ?) OR (`ip` LIKE ? OR `ipx` LIKE ? OR `port` LIKE ? OR `protocol` LIKE ? OR `url` LIKE ? OR `location` LIKE ? OR `title` LIKE ? OR `Host` LIKE ?) OR `Host` LIKE ?
PARAMS: [%127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %1234% %1234% %1234% %1234% %1234% %1234% %1234% %1234% %1%]
```