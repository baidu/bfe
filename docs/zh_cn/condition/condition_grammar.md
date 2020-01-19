# 条件语法


## 基本概念: 原语、表达式和变量

- **条件原语**（Condition Primitive）

  - 最基本的条件判断单元，定义了比较的原语

  ``` 
  req_host_in("www.bfe-networks.com|bfe-networks.com") # host是两个域名之一
  ```
  
- **条件表达式**（Condition Expression）：基于条件原语的“与/或/非组合”

  ```
  req_host_in("bfe-networks.com") && req_method_in("GET") # 域名是bfe-networks.com且方法是GET
  ```
  
- **条件变量**（Condition Variable）

  - 可以将**条件表达式**赋值给一个变量，这个变量被定义为条件变量

  ```
  bfe_host = req_host_in("bfe-networks.com") # 将条件表达式赋值给变量bfe_host
  ```
  
- **高级条件表达式**（Advanced Condition Expression）：基于条件原语和条件变量的“与/或/非组合”

  - 在高级条件表达式中，条件变量以$前缀作为标示

  ```
  $news_host && req_method_in("GET") # 符合变量news_host且方法是GET
  ```


## 条件原语的语法

- **条件原语**是最基本的条件判断单元，形式为：

​           **func_name( params )**

- 条件原语类似于函数定义: func_name是条件原语名称; params是参数，可能是1个或多个
- 条件原语的返回值都是bool类型
- BFE支持一系列预定义的内置[条件原语](condition_primitive_index.md)
  - 例: req_method_in("GET|POST") ，表示判断http请求的方法是不是GET或POST


## 条件表达式的语法

条件表达式的语法定义如下：

- 与c语言中的&&、||、！有一致的优先级和结合律

- 语法描述

  ```
  Condition Expression(CE) -> 
                     CE && CE
                   | CE || CE
                   | ( CE )
                   | ! CE
                   | Condition Primitive
  ```
  

## 高级条件表达式的语法

高级条件表达式的语法定义如下：

- 与c语言中的&&、||、！有一致的优先级和结合律

- 语法描述

  ```
  Advanced Condition Expression(ACE) -> 
                     ACE && ACE
                   | ACE || ACE
                   | ( ACE)
                   | ! ACE
                   | Condition Primitive
                   | Condition Variable
  ```
  
