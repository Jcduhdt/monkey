# Monkey Language
development by go

# document
## ast
- abstract syntax code 
- 抽象语法树
## evaluator
- 求值器
- 将表达式进行取值
## lexer
- 词法分析器 
- 将源代码转换成词法单元
## parser
- 语法分析器 
- 采用递归下降语法分析
### 普拉特语法分析器
  - 主要思想
    - 将解析函数与词法单元类型相关联。每当遇到某个词法单元类型时，都会调用相关联的解析函数来解析对应的表达式，最后返回生成的ast节点。
    - 每个词法单元类型最多可以关联两个解析函数，取决于词法单元位置，前缀or中缀
## repl
- Read-Eval-Print Loop
## token
- 词法单元