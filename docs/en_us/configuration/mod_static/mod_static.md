# Introduction 

mod_static serves static files.

# Module configuration

## Description
conf/mod_static/mod_static.conf

| Config Item | Description                             |
| ----------- | --------------------------------------- |
| Basic.DataPath | String<br>path of rule configuraiton |

## Example
```
[basic]
DataPath = mod_static/static_rule.data
```

# Rule configuration

## Description
conf/mod_static/static_rule.data

| Config Item | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| Version     | String<br>Verson of config file |
| Config      | Struct<br>Static rules for each product |
| Config{k}   | String<br>Product name |
| Config{v}   | Object<br>A ordered list of static rules |
| Config{v}[] | Object<br>A static rule |
| Config{v}[].Cond | String<br>Condition expression, See [Condition](../../condition/condition_grammar.md) |
| Config{v}[].Action | Object<br>A static action |
| Config{v}[].Action.Cmd | String<br>Name of static action |
| Config{v}[].Action.Params | Object<br>Parameters of static action |

## Actions
| Action                    | Description                        |
| ------------------------- | ---------------------------------- |
| BROWSE                    | Serve static files. <br>The first parameter is the location of root directory.<br> The second parameter is the name of default file.|

## Example
```
{
    "Config": {
        "example_product": [
            {
                "Cond": "req_host_in(\"www.example.org\")",
                "Action": {
                    "Cmd": "BROWSE",
                    "Params": [
                        "./",
                        "index.html"
                    ]
                }
            }
        ]
    },
    "Version": "20190101000000"
}
```
