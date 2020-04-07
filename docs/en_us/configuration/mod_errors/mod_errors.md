# Introduction

mod_errors replaces error responses by specified rules.

# Module configuration
## Description
conf/mod_errors/mod_errors.conf

| Config Item          | Description                                 |
| ---------------------| ------------------------------------------- |
| Basic.DataPath       | String<br>Path fo rule configuration |
| Log.OpenDebug        | Boolean<br>Whether enable debug logs<br>Default False |

## Example
```
[basic]
DataPath = mod_errors/errors_rule.data
```

# Rule configuration
## Description 

| Config Item | Description                                                |
| ----------- | ---------------------------------------------------------- |
| Version | String<br>Version of config file |
| Config | Object<br>Error rules for each product |
| Config{k} | String<br>Product name |
| Config{v} | Object<br> A list of error rules |
| Config{v}[] | Object<br>A error rule |
| Config{v}[].Cond | String<br>Condition expressio, See [Condition](../../condition/condition_grammar.md) |
| Config{v}[].Actions | Object<br>Action |
| Config{v}[].Actions.Cmd | String<br>Name of Action |
| Config{v}[].Actions.Params | Object<br>Parameters of Action |
| Config{v}[].Actions.Params[] | String<br>A Parameter |

## Module actions
| Action   | Description            |
| -------- | ---------------------- |
| RETURN   | Return response generated from specified static html |
| REDIRECT | Redirect to specified location |

## Example
```
{
    "Version": "20190101000000",
    "Config": {
        "example_product": [
            {
                "Cond": "res_code_in(\"404\")",
                "Actions": [
                    {
                        "Cmd": "RETURN",
                        "Params": [
                            "200", "text/html", "../conf/mod_errors/404.html"
                        ]
                    }
                ]
            },
            {
                "Cond": "res_code_in(\"500\")",
                "Actions": [
                    {
                        "Cmd": "REDIRECT",
                        "Params": [
                            "http://example.org/error.html"
                        ]
                    }
                ]
            }
        ]
    }
}
```

