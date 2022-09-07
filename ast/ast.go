package ast

import (
	"monkey/token"
)

// Node AST中每个节点都必须实现Node接口
type Node interface {
	// 该方法返回与其关联的词法单元的字面量
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

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

type LetStatement struct {
	Token token.Token // token.Let 词法单元
	Name  *Identifier // 标识符，为了减少AST中各种类型节点的数量，复用该节点
	Value Expression  // 产生值的表达式
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT 词法单元
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
