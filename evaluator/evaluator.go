package evaluator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/govel-framework/lamb/ast"
	"github.com/govel-framework/lamb/internal"
	"github.com/govel-framework/lamb/object"
	"github.com/govel-framework/lamb/token"
)

func Eval(node ast.Node, env *object.Environment) interface{} {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return node.Value

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right, node.Token)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right, node.Token)

	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.VarStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.CallExpression:
		function := Eval(node.Function, env)

		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {

			return args[0]
		}

		return applyFunction(function, args, node.Token)

	case *ast.StringLiteral:
		if !node.Closed {
			return newError(node.Token, "unclosed string literal")
		}

		return node.Value

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)

		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return elements

	case *ast.IndexExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)

		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index, node.Token)

	case *ast.MapLiteral:
		return evalMapLiteral(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.ExtendsStatement:
		return evalExtendsStatement(node, env)

	case *ast.SectionStatement:
		return evalSectionStatement(node, env)

	case *ast.DefineStatement:
		return evalDefineStatement(node, env)

	case *ast.DotExpression:
		return evalDotExpression(node, env)

	case *ast.IncludeStatement:
		return evalIncludeStatement(node, env)

	case *ast.HtmlLiteral:
		return node.Value
	}

	return nil
}

func evalStatements(stmts []ast.Statement, env *object.Environment) interface{} {
	// save the result as a string
	var result string

	for _, statement := range stmts {
		res := Eval(statement, env)

		if isError(res) {
			return res
		}

		if res != nil {
			result += fmt.Sprintf("%v", res)
		}

	}

	// return the output of the statements as an object.String
	return result
}

func nativeBoolToBooleanObject(input bool) bool {
	if input {
		return true
	}

	return false
}

func evalPrefixExpression(operator string, right interface{}, t token.Token) interface{} {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)

	case "-":
		return evalMinusPrefixOperatorExpression(right, t)

	default:
		return newError(t, "unknown operator: %s%T", operator, right)
	}
}

func evalBangOperatorExpression(right interface{}) interface{} {
	switch right {
	case true:
		return false
	case false:
		return true
	case nil:
		return true
	default:
		return false
	}
}

func evalMinusPrefixOperatorExpression(right interface{}, t token.Token) interface{} {
	if fmt.Sprintf("%T", right) != "int" {
		return newError(t, "unknown operator: -%T", right)
	}

	value := right.(int)
	return -value
}

func evalInfixExpression(operator string, left, right interface{}, t token.Token) interface{} {
	leftType := fmt.Sprintf("%T", left)
	rightType := fmt.Sprintf("%T", right)

	leftNumber, isLeftNumber := isNumber(left)
	rightNumber, isRightNumber := isNumber(right)

	switch {
	case isLeftNumber && isRightNumber:
		return evalIntegerInfixExpression(operator, leftNumber, rightNumber, t)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case operator == "and":
		leftValue := reflect.ValueOf(left)
		rightValue := reflect.ValueOf(right)

		if (leftValue.Kind() == reflect.Ptr && leftValue.IsNil()) || !leftValue.IsValid() {
			return false
		} else if boolValue, isBool := right.(bool); isBool && !boolValue {
			return false
		}

		if (rightValue.Kind() == reflect.Ptr && rightValue.IsNil()) || !rightValue.IsValid() {
			return false
		} else if boolValue, isBool := right.(bool); isBool && !boolValue {
			return false
		}

		return true

	case leftType == "string" && rightType == "string":
		return evalStringInfixExpression(operator, left, right, t)

	case leftType != rightType:
		return newError(t, "type mismatch: %s %s %s", leftType, operator, rightType)

	default:
		return newError(t, "unknown operator: %s %s %s", leftType, operator, rightType)
	}
}

func evalIntegerInfixExpression(operator string, left, right interface{}, t token.Token) interface{} {
	leftVal := left.(int)

	rightVal := right.(int)

	switch operator {
	case "+":
		return leftVal + rightVal

	case "-":
		return leftVal - rightVal

	case "*":
		return leftVal * rightVal

	case "/":
		return leftVal / rightVal

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)

	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)

	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)

	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return newError(t, "unknown operator: %T %s %T", left, operator, right)
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) interface{} {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)

	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)

	} else {
		return nil
	}

}

func isTruthy(obj interface{}) bool {
	switch obj {

	case nil:
		return false

	case true:
		return true

	case false:
		return false

	default:
		return true
	}
}

func newError(t token.Token, format string, a ...interface{}) error {
	err := fmt.Sprintf("%d: %d: ", t.Line, t.Col)

	return fmt.Errorf(err+format, a...)
}

func evalProgram(program *ast.Program, env *object.Environment) interface{} {
	var result string

	for _, statement := range program.Statements {
		r := Eval(statement, env)

		if isError(r) {
			return fmt.Sprintf("%s: %v", env.FileName, r)
		}

		if r != nil {
			result += fmt.Sprintf("%v", r)
		}
	}

	if env.InExtends {
		// eval the file and create the new environment
		newEnv := object.CopyEnvironment(env)
		newEnv.IsExtends = true

		var out bytes.Buffer

		err := internal.LoadFile(env.ExtendsFrom.From, nil, &out, Eval, *newEnv)

		result = out.String()

		// check if any error has occured
		if err != nil {
			return errors.New(err.Error())
		}

		// check if any section is ununsed
		for _, section := range env.ExtendsFrom.Sections {
			return newError(section.Token, "section %s does not exist", section.Name)
		}

	}

	return result
}

func isError(obj interface{}) bool {
	if obj != nil {
		_, is := obj.(error)

		return is
	}

	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) interface{} {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := Builtins[node.Value]; ok {
		return builtin
	}

	return newError(node.Token, "identifier not found: %s", node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []interface{} {
	var result []interface{}

	for _, e := range exps {
		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []interface{}{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn interface{}, args []interface{}, t token.Token) interface{} {
	switch fn := fn.(type) {

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError(t, "not a function: %T", fn)
	}
}

func evalStringInfixExpression(operator string, left, right interface{}, t token.Token) interface{} {
	if operator != "+" {
		return newError(t, "unknown operator: %T %s %T", left, operator, right)
	}

	leftVal := left.(string)

	rightVal := right.(string)

	return leftVal + rightVal
}

func evalIndexExpression(left, index interface{}, t token.Token) interface{} {
	leftType := reflect.ValueOf(left).Kind()
	indexType := reflect.ValueOf(index).Kind()

	switch {
	case (leftType == reflect.Slice || leftType == reflect.Array) && indexType == reflect.Int:
		return evalArrayIndexExpression(left, index)

	case leftType == reflect.Map:
		return evalMapIndexExpression(left, index)

	default:
		return newError(t, "index operator not supported: %s", leftType.String())
	}
}

func evalArrayIndexExpression(array, index interface{}) interface{} {
	arrayValue := reflect.ValueOf(array)

	id := index.(int)

	max := arrayValue.Len() - 1

	if id < 0 || id > max {
		return nil
	}

	return arrayValue.Index(id).Interface()
}

func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) interface{} {
	pairs := make(map[interface{}]interface{})

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)

		if isError(key) {
			return key
		}

		value := Eval(valueNode, env)

		if isError(value) {
			return value
		}

		pairs[key] = value

	}

	return pairs
}

func evalMapIndexExpression(m, index interface{}) interface{} {
	mapValue := reflect.ValueOf(m)

	value := mapValue.MapIndex(reflect.ValueOf(index))

	if !value.IsValid() {
		return nil
	}

	return value.Interface()
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) interface{} {
	value := fe.Value
	key := fe.Key

	in := Eval(fe.In, env)

	if isError(in) {
		return in
	}

	// iterate
	var out string

	valueOf := reflect.ValueOf(in)

	switch valueOf.Kind() {

	case reflect.Map:

		for _, elem := range valueOf.MapKeys() {

			// set the new values
			env.Set(value, elem.Interface())

			if key != "" {
				env.Set(key, elem.Interface())
			}

			res := Eval(fe.Block, env)

			if isError(res) {
				return res
			}

			out += res.(string)
		}

	case reflect.Array, reflect.Slice:
		len := valueOf.Len()

		for i := 0; i < len; i++ {
			elem := valueOf.Index(i).Interface()

			// set the new values
			env.Set(value, elem)

			if key != "" {
				env.Set(key, i)
			}

			res := Eval(fe.Block, env)

			if isError(res) {
				return res
			}

			out += res.(string)
		}

	default:
		return newError(fe.Token, "%T is not iterable", in)
	}

	// delete the vars
	env.Delete(value)

	if key != "" {
		env.Delete(key)
	}

	return out
}

func evalExtendsStatement(node *ast.ExtendsStatement, env *object.Environment) interface{} {
	if env.InExtends || env.IsExtends {
		return newError(node.Token, "nested extends are not allowed")
	}

	env.InExtends = true
	env.ExtendsFrom.From = node.From

	return nil
}

func evalSectionStatement(node *ast.SectionStatement, env *object.Environment) interface{} {
	if !env.InExtends {
		return newError(node.Token, "section statement is only allowed in extends")
	}

	if env.IsExtends {
		return newError(node.Token, "section statement is only allowed with extends")
	}

	if env.InSection {
		return newError(node.Token, "section statement is not allowed in a section")
	}

	// save the section
	env.ExtendsFrom.Sections[node.Name] = object.SectionContent{
		Content: Eval(node.Block, env),
		Name:    node.Name,
		Token:   node.Token,
	}

	return nil
}

func evalDefineStatement(node *ast.DefineStatement, env *object.Environment) interface{} {
	var content interface{}

	if env.InDefine {
		return newError(node.Token, "nested defines are not allowed")
	}

	// check if the section exists
	if section, ok := env.ExtendsFrom.Sections[node.Name]; ok {
		content = section.Content

		// delete the section
		delete(env.ExtendsFrom.Sections, node.Name)

	} else {
		content = Eval(node.Content, env)
	}

	return content
}

func evalDotExpression(node *ast.DotExpression, env *object.Environment) interface{} {
	var result interface{}

	left := Eval(&node.Left, env)

	if isError(left) {
		return left
	}

	leftValue := reflect.ValueOf(left)
	leftType := reflect.ValueOf(left).Kind()

	if leftType == reflect.Ptr {
		leftValue = leftValue.Elem()

		leftType = leftValue.Kind()
	}

	if leftType != reflect.Struct {
		return newError(node.Token, "left side of dot expression must be a struct, got=%s", leftType)
	}

	leftStruct := reflect.TypeOf(leftValue.Interface())

	// check if the field (node.Right) exists
	if _, ok := leftStruct.FieldByName(node.Right.Value); ok {

		result = leftValue.FieldByName(node.Right.Value).Interface()

	} else {
		return newError(node.Token, "field %s does not exist in struct %s", node.Right.Value, node.Left.Value)
	}

	return result
}

func isNumber(num interface{}) (int, bool) {
	if num == nil {
		return 0, false
	}

	is := reflect.TypeOf(num).Kind() == reflect.Int ||
		reflect.TypeOf(num).Kind() == reflect.Int8 ||
		reflect.TypeOf(num).Kind() == reflect.Int16 ||
		reflect.TypeOf(num).Kind() == reflect.Int32 ||
		reflect.TypeOf(num).Kind() == reflect.Int64

	if !is {
		return 0, false
	}

	return int(reflect.ValueOf(num).Int()), true
}

func evalIncludeStatement(node *ast.IncludeStatement, env *object.Environment) interface{} {
	newEnv := object.NewEnvironment()

	if node.Vars != nil {
		vars, isMap := node.Vars.(*ast.MapLiteral)

		if !isMap {
			return newError(node.Token, "vars in include must be a map, got=%s", node.Vars.TokenLiteral())
		}

		for key, value := range vars.Pairs {
			newEnv.Set(key.String(), Eval(value, env))
		}
	}

	var out bytes.Buffer

	err := internal.LoadFile(node.File, nil, &out, Eval, *newEnv)

	result := out.String()

	// check if any error has occured
	if err != nil {
		return errors.New(err.Error())
	}

	return result
}
