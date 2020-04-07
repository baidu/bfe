# Condition Grammar

## Basic Concepts

- **Condition Primitive**

  - Basic conditional judgment unit, which defines the primitive of comparison

    ``` 
    req_host_in("www.bfe-networks.com|bfe-networks.com")  # host is one of the configured domains
    ```

- **Condition Expression**

  - Expression using "and/or/not" to connect condition primitive

    ```
    req_host_in("bfe-networks.com") && req_method_in("GET") # domain is bfe-networks.com and HTTP method is "GET"
    ```

- **Condition Variable**

  - Variable that is defined by **Condition Expression**

    ```
    bfe_host = req_host_in("bfe-networks.com")  # variable bfe_host is defined by condition expression 
    ```

- **Advanced Condition Expression**

  - Expression using "and/or/not" to connect condition primitive and condition variable

  - In advanced condition expression, condition variable is identified by  **"$" prefix**

    ```
    $news_host && req_method_in("GET") # match condition variable and HTTP method is "GET"
    ```


## Condition Primitive Grammar

- Basic conditional judgment unit, format is shown as follows:

​           **FuncName( params )**

- Condition primitive like function definition: FuncName is name of condition primitive; params are input parameters
- Return value type of Condition Primitive is bool
- Note: All builtin [condition primitives](condition_primitive_index.md)


## Condition Expression Grammar

Condition Expression grammar is defined as follows:

- Priority and combination rule of "&&/||/!" is same as them in C language

- Expression description

  ```
  Condition Expression(CE) -> 
                     CE && CE
                   | CE || CE
                   | ( CE )
                   | ! CE
                   | Condition Primitive
  ```
  
  

## Advanced Condition Expression Grammar

Advanced Condition Expression grammar is defined as follows:

- Priority and combination rule of "&&/||/!" is same as them in C language

- Expression description

  ```
  Advanced Condition Expression(ACE) -> 
                     ACE && ACE
                   | ACE || ACE
                   | ( ACE)
                   | ! ACE
                   | Condition Primitive
                   | Condition Variable
  ```
  
  
