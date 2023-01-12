## 简介

> 逻辑查询解析器

## TODO

- [ ] 为keyword添加格式检测，如ip="1.1.1.1"(正确), ip="1"(错误)

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

## TEST

### GORM

```bash
 INPUT: IP="127.0.0.1"||127.0.0.1 || 1234 || HOST=1
FORMAT: ip="127.0.0.1" || "127.0.0.1" || "1234" || Host="1"
   SQL: `ip` LIKE ? OR (`ip` LIKE ? OR `ipx` LIKE ? OR `port` LIKE ? OR `protocol` LIKE ? OR `url` LIKE ? OR `location` LIKE ? OR `title` LIKE ? OR `Host` LIKE ?) OR (`ip` LIKE ? OR `ipx` LIKE ? OR `port` LIKE ? OR `protocol` LIKE ? OR `url` LIKE ? OR `location` LIKE ? OR `title` LIKE ? OR `Host` LIKE ?) OR `Host` LIKE ?
PARAMS: [%127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %127.0.0.1% %1234% %1234% %1234% %1234% %1234% %1234% %1234% %1234% %1%]
```

### bleve

```bash
FORMAT: (ip="1" || ip="2") && protocol=="https" && "127.0.0.1" && ((title*="1") || title!*="2") && ip="1" && (((title="1" || title="2")) && ip="2" || (ip="1" || ((ip="10") || ip="20") && title="ccccc")) 
```

```json
{
    "query": {
        "conjuncts": [
            {
                "conjuncts": [
                    {
                        "conjuncts": [
                            {
                                "conjuncts": [
                                    {
                                        "conjuncts": [
                                            {
                                                "disjuncts": [
                                                    {
                                                        "match": "1",
                                                        "field": "ip",
                                                        "prefix_length": 0,
                                                        "fuzziness": 0
                                                    },
                                                    {
                                                        "match": "2",
                                                        "field": "ip",
                                                        "prefix_length": 0,
                                                        "fuzziness": 0
                                                    }
                                                ],
                                                "min": 0
                                            },
                                            {
                                                "term": "https",
                                                "field": "protocol"
                                            }
                                        ]
                                    },
                                    {
                                        "match": "127.0.0.1",
                                        "prefix_length": 0,
                                        "fuzziness": 0
                                    }
                                ]
                            },
                            {
                                "disjuncts": [
                                    {
                                        "wildcard": "1",
                                        "field": "title"
                                    },
                                    {
                                        "must_not": {
                                            "disjuncts": [
                                                {
                                                    "wildcard": "2",
                                                    "field": "title"
                                                }
                                            ],
                                            "min": 0
                                        }
                                    }
                                ],
                                "min": 0
                            }
                        ]
                    },
                    {
                        "match": "1",
                        "field": "ip",
                        "prefix_length": 0,
                        "fuzziness": 0
                    }
                ]
            },
            {
                "disjuncts": [
                    {
                        "conjuncts": [
                            {
                                "disjuncts": [
                                    {
                                        "match": "1",
                                        "field": "title",
                                        "prefix_length": 0,
                                        "fuzziness": 0
                                    },
                                    {
                                        "match": "2",
                                        "field": "title",
                                        "prefix_length": 0,
                                        "fuzziness": 0
                                    }
                                ],
                                "min": 0
                            },
                            {
                                "match": "2",
                                "field": "ip",
                                "prefix_length": 0,
                                "fuzziness": 0
                            }
                        ]
                    },
                    {
                        "conjuncts": [
                            {
                                "disjuncts": [
                                    {
                                        "match": "1",
                                        "field": "ip",
                                        "prefix_length": 0,
                                        "fuzziness": 0
                                    },
                                    {
                                        "disjuncts": [
                                            {
                                                "match": "10",
                                                "field": "ip",
                                                "prefix_length": 0,
                                                "fuzziness": 0
                                            },
                                            {
                                                "match": "20",
                                                "field": "ip",
                                                "prefix_length": 0,
                                                "fuzziness": 0
                                            }
                                        ],
                                        "min": 0
                                    }
                                ],
                                "min": 0
                            },
                            {
                                "match": "ccccc",
                                "field": "title",
                                "prefix_length": 0,
                                "fuzziness": 0
                            }
                        ]
                    }
                ],
                "min": 0
            }
        ]
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