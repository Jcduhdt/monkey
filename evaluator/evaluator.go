package evaluator

import (
	"fmt"

	"monkey/ast"
	"monkey/object"
)

var (
	// null只有一种
	NULL = &object.Null{}
	// bool只有两种可能，所以不必每次都创建实例。使用以下两个引用代替每次新实例
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch nodeT := node.(type) {
	case *ast.Program:
		return evalProgram(nodeT, env)
	case *ast.ExpressionStatement:
		return Eval(nodeT.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: nodeT.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(nodeT.Value)
	case *ast.PrefixExpression:
		right := Eval(nodeT.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(nodeT.Operator, right)
	case *ast.InfixExpression:
		left := Eval(nodeT.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(nodeT.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(nodeT.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(nodeT, env)
	case *ast.IfExpression:
		return evalIfExpression(nodeT, env)
	case *ast.ReturnStatement:
		val := Eval(nodeT.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(nodeT.Value, env)
		if isError(val) {
			return val
		}
		// 哈希映射
		env.Set(nodeT.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(nodeT, env)
	case *ast.FunctionLiteral:
		params := nodeT.Parameters
		body := nodeT.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(nodeT.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(nodeT.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: nodeT.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(nodeT.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(nodeT.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(nodeT.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch rt := result.(type) {
		case *object.ReturnValue:
			return rt.Value
		case *object.Error:
			return rt
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	}
	return newError("unknown operator: %s%s", operator, right.Type())
}

// 支持!5这种为false
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	// 如果当前环境没有发现对应的标识符，则在内置函数环境中查找
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fnT := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fnT, args)
		evaluated := Eval(fnT.Body, extendedEnv)
		return unWarpReturnValue(evaluated)
	case *object.Builtin:
		return fnT.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

// 扩展的是定义函数时的环境，而不是当前环境。闭包得以实现
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		// 为什么不是env.Set(param.Value,param)呢
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

// 解包是为了避免return向上冒泡，使外层函数停止取值
func unWarpReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}
