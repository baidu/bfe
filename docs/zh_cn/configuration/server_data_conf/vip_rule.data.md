# 简介

vip_rule.data配置文件记录产品线的VIP列表。

# 配置

| 配置项  | 类型   | 描述                                                         |
| ------- | ------ | ------------------------------------------------------------ |
| Version | String | 配置文件版本                                                 |
| Vips    | Map&lt;String, Array&lt;String&gt;&gt; | 产品线的VIP列表，key是产品线名称，value是VIP列表 |

# 示例

```
{
    "Version": "20190101000000",
    "Vips": {
        "example_product": [
            "111.111.111.111"
        ] 
    }
}
```
