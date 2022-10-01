package evaluator

import (
	"fmt"

	"monkey/ast"
	"monkey/object"
	"monkey/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

// 对unquote的内容进行解析
func evalUnquoteCalls(quote ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quote, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}
		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return callExpression.Function.TokenLiteral() == "unquote"
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch objT := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", objT.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: objT.Value}
	case *object.Boolean:
		var t token.Token
		if objT.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: objT.Value}
	case *object.Quote:
		return objT.Node
	default:
		return nil
	}
}
