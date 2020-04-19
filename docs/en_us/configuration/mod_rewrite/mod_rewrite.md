# Introduction 

Modify URI of HTTP request based on defined rules.

# Module configuration

## Description
conf/mod_rewrite/mod_rewrite.conf

| Config Item | Description                             |
| ----------- | --------------------------------------- |
| Basic.DataPath | String<br>path of rule configuraiton |

## Example

```
[basic]
DataPath = mod_rewrite/rewrite.data
```

# Rule configuration

## Description
conf/mod_rewrite/rewrite.data

| Config Item | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| Version     | String<br>Verson of config file |
| Config      | Struct<br>Rewrite rules for each product |
| Config{k}   | String<br>Product name |
| Config{v}   | Object<br>A ordered list of rewrite rules |
| Config{v}[] | Object<br>A rewrite rule |
| Config{v}[].Cond | String<br>Condition expression, See [Condition](../../condition/condition_grammar.md) |
| Config{v}[].Actions | Object<br>A ordered list of rewrite actions |
| Config{v}[].Actions[] | Object<br>A rewrite action |
| Config{v}[].Actions[].Cmd | Object<br>Name of rewrite action |
| Config{v}[].Actions[].Params | Object<br>Parameters of rewrite action |
| Config{v}[].Last | Integer<br>If true, stop to check the remaining rules |


## Actions
| Action                    | Description                              |
| ------------------------- | ---------------------------------------- |
| HOST_SET                  | Set host to specified value              |
| HOST_SET_FROM_PATH_PREFIX | Set host to specified path prefix        |
| PATH_SET                  | Set path to specified value              |
| PATH_PREFIX_ADD           | Add prefix to orignal path               |
| PATH_PREFIX_TRIM          | Trim prefix from orignal path            |
| QUERY_ADD                 | Add query                                |
| QUERY_DEL                 | Delete query                             |
| QUERY_DEL_ALL_EXCEPT      | Del all queries except specified queries |
| QUERY_RENAME              | Rename query                             |
  
## Example

```
{
    "Version": "20190101000000",
    "Config": {
        "example_product": [
            {
                "Cond": "req_path_prefix_in(\"/rewrite\", false)",
                "Actions": [
                    {
                        "Cmd": "PATH_PREFIX_ADD",
                        "Params": [
                            "/bfe/"
                        ]
                    }
                ],
                "Last": true
            }
        ]
    }
}
```
  
