package ast

import (
	"bytes"
	"strings"

	"monkey/token"
)

// Node AST中每个节点都必须实现Node接口
type Node interface {
	// 该方法返回与其关联的词法单元的字面量
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

// 表达式的具体内容
type Expression interface {
	Node
	expressionNode()
}

// Program 语法分析器生成的每个AST的根节点
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// let语句解析
type LetStatement struct {
	Token token.Token // token.Let 词法单元
	Name  *Identifier // 标识符，为了减少AST中各种类型节点的数量，复用该节点
	Value Expression  // 产生值的表达式
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token // token.IDENT 词法单元
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string {
	return i.Value
}

// return语句解析
type ReturnStatement struct {
	Token       token.Token // return词法单元
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// 表达式语句
type ExpressionStatement struct {
	Token      token.Token // 该表达式中的第一个词法单元
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *IntegerLiteral) String() string { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

func (b *Boolean) String() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {}

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression // 标识符或函数字面量
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

func (sl *StringLiteral) String() string { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// 语法分析阶段，所有表达式都应该可以用做哈希字面量中的键和值
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, val := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+val.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type modifierFunc func(Node) Node

// 只修改了子节点，没有修改父节点会导致String()输出不一致。
// 实时创建新的词法单元，会丢失其来源信息
func Modify(node Node, modifier modifierFunc) Node {
	switch nodeT := node.(type) {
	case *Program:
		for i, statement := range nodeT.Statements {
			nodeT.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		nodeT.Expression, _ = Modify(nodeT.Expression, modifier).(Expression)
	case *InfixExpression:
		nodeT.Left, _ = Modify(nodeT.Left, modifier).(Expression)
		nodeT.Right, _ = Modify(nodeT.Right, modifier).(Expression)
	case *PrefixExpression:
		nodeT.Right, _ = Modify(nodeT.Right, modifier).(Expression)
	case *IndexExpression:
		nodeT.Left, _ = Modify(nodeT.Left, modifier).(Expression)
		nodeT.Index, _ = Modify(nodeT.Index, modifier).(Expression)
	case *IfExpression:
		nodeT.Condition, _ = Modify(nodeT.Condition, modifier).(Expression)
		nodeT.Consequence, _ = Modify(nodeT.Consequence, modifier).(*BlockStatement)
		if nodeT.Alternative != nil {
			nodeT.Alternative, _ = Modify(nodeT.Alternative, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		for i, _ := range nodeT.Statements {
			nodeT.Statements[i], _ = Modify(nodeT.Statements[i], modifier).(Statement)
		}
	case *ReturnStatement:
		nodeT.ReturnValue, _ = Modify(nodeT.ReturnValue, modifier).(Expression)
	case *LetStatement:
		nodeT.Value, _ = Modify(nodeT.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i, _ := range nodeT.Parameters {
			nodeT.Parameters[i], _ = Modify(nodeT.Parameters[i], modifier).(*Identifier)
		}
		nodeT.Body = Modify(nodeT.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i, _ := range nodeT.Elements {
			nodeT.Elements[i], _ = Modify(nodeT.Elements[i], modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, val := range nodeT.Pairs {
			newKey, _ := Modify(key, modifier).(Expression)
			newVal, _ := Modify(val, modifier).(Expression)
			newPairs[newKey] = newVal
		}
		nodeT.Pairs = newPairs
	}
	return modifier(node)
}

type MacroLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ml *MacroLiteral) expressionNode() {}

func (ml *MacroLiteral) TokenLiteral() string { return ml.Token.Literal }

func (ml *MacroLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range ml.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(ml.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(ml.Body.String())

	return out.String()
}
