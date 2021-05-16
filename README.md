# README

配置示例：
```
{
    @default("xxxxxxxxxx")
    title

    @css(".items")
    @return
    items [
        {
            @css(".title")
            @text()
            title
        
            @css(".link")
            @attr(".attr")
            @text()
            link

            @link(".tags")
            tags [
                @trim()
                @float()
                value
            ]
        }
    ]
}
```

解析后语法：
```json
[
    {
        "name": "items",
        "pipelines": [
            {
                "method": "css",
                "arguments": [".items"]
            }
        ],
        "children": [
            {
                "name": "title",
                "pipelines": [
                    {
                        "method": "css",
                        "arguments": [".title"]
                    }
                ]
            },
            {
                "name": "link",
                "pipelines": [
                    {
                        "method": "css",
                        "arguments": [".link"]
                    }
                ]
            },
            {
                "name": "tags",
                "pipelines": [
                    {
                        "method": "css",
                        "arguments": [".tags"]
                    }
                ],
                "children": [
                    {
                        "name": "value",
                        "pipelines": [
                            {
                                "method": "trim"
                            }                                }
                        ]
                    }
                ]
            }
        ]
    }
]
```


## 支持默认值
支持默认值，比如默认字符串、当前时间

## 参考其他定义
比如protobuf、ts、go等，是不是可以复用其语法解析规则，简化实现与用户学习成本