package evaluator

import (
	"fmt"
	"reflect"

	"github.com/govel-framework/govel"

	"github.com/govel-framework/lamb/object"
)

func builtInError(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// Builtins is a map of builtin functions.
//
// DO NOT USE THIS MAP DIRECTLY as it is for private use only.
var Builtins = map[string]*object.Builtin{
	"len": {
		Fn: lenBuiltIn,
	},
	"type": {
		Fn: typeBuiltIn,
	},
	"map_key_exists": {
		Fn: mapKeyExists,
	},
	"range": {
		Fn: rangeBuiltIn,
	},
	"route": {
		Fn: routeBuiltIn,
	},
	"config": {
		Fn: configBuiltIn,
	},
	"asset": {
		Fn: assetBuiltIn,
	},
}

func lenBuiltIn(args ...interface{}) interface{} {
	if len(args) != 1 {
		return builtInError("wrong number of arguments in len. got=%d, want=1", len(args))
	}

	valueOf := reflect.ValueOf(args[0])

	switch valueOf.Kind() {

	case reflect.Slice, reflect.Array:
		return int64(valueOf.Len())

	case reflect.String:
		return int64(len(args[0].(string)))

	default:
		return builtInError("argument to `len` not supported, got %T", args[0])
	}

}

func typeBuiltIn(args ...interface{}) interface{} {
	if len(args) != 1 {
		return builtInError("wrong number of arguments in type. got=%d, want=1", len(args))
	}

	arg := args[0]

	return fmt.Sprintf("%T", arg)
}

func mapKeyExists(args ...interface{}) interface{} {
	if len(args) != 2 {
		return builtInError("wrong number of arguments in map_key_exists. got=%d, want=2", len(args))
	}

	m := args[0]
	k := args[1]

	// validate the type of m

	if m == nil {
		return false
	}

	mValue := reflect.ValueOf(m)

	if mValue.Kind() != reflect.Map {
		return builtInError("argument to `map_key_exists` not supported, got %T, want=map", m)
	}

	// check if the key exists
	value := mValue.MapIndex(reflect.ValueOf(k))

	if !value.IsValid() {
		return false
	}

	return true
}

func rangeBuiltIn(args ...interface{}) interface{} {
	if len(args) != 2 {
		return builtInError("wrong number of arguments in range. got=%d, want=2", len(args))
	}

	start := args[0]

	end := args[1]

	if reflect.TypeOf(start).Kind() != reflect.Int || reflect.TypeOf(end).Kind() != reflect.Int {
		return builtInError("argument to `range` not supported, got %T, want=int", args[0])
	}

	startInt := start.(int)
	endInt := end.(int)

	result := []int{}

	for startInt <= endInt {
		result = append(result, startInt)
		startInt++
	}

	return result
}

func routeBuiltIn(args ...interface{}) interface{} {
	routeArgs := make(map[interface{}]string)

	// validate the number of arguments, min 1: string, max 2: map[string]string
	if len(args) < 1 {
		return builtInError("wrong number of arguments in route. got=%d, want=1", len(args))
	}

	route := args[0]

	if fmt.Sprintf("%T", route) != "string" {
		return builtInError("argument to `route` not supported, got %T, want=string", route)
	}

	if len(args) > 2 {
		return builtInError("wrong number of arguments in route. got=%d, want=2", len(args))
	}

	if len(args) == 2 {
		rArgs, isMap := args[1].(map[interface{}]interface{})

		if !isMap {
			return builtInError("argument to `route` not supported, got %T, want=map", rArgs)
		}

		for key, value := range rArgs {
			keyType := reflect.TypeOf(key).Kind()
			valueType := reflect.TypeOf(value).Kind()

			if keyType != reflect.String && keyType != reflect.Int {
				return builtInError("argument to `route` not supported, all elements of map must be strings or integers. got=%s", keyType)
			}

			if valueType != reflect.String && valueType != reflect.Int {
				return builtInError("argument to `route` not supported, all elements of map must be strings or integers. got=%s", valueType)
			}

			routeArgs[key] = fmt.Sprintf("%v", value)
		}
	}

	// convert routeArgs (map[interface{}]string) to (map[string]string)
	routeArgsString := make(map[string]string)

	for key, value := range routeArgs {
		routeArgsString[fmt.Sprintf("%v", key)] = value
	}

	url := govel.Route(route.(string), routeArgsString)

	if url == "" {
		panic(fmt.Sprintf("Route %s not found", route))
	}

	return url
}

func configBuiltIn(args ...interface{}) interface{} {
	if len(args) != 1 {
		return builtInError("wrong number of arguments in config. got=%d, want=1", len(args))
	}

	arg := args[0]
	argType := reflect.TypeOf(arg).Kind()

	if argType != reflect.String {
		return builtInError("argument to `config` not supported, got %s, want=string", argType.String())
	}

	// split the string
	key := arg.(string)

	exists, value := lookForConfigKeys(govel.GetKeyFromYAML("").(map[interface{}]interface{}), key)

	if !exists {
		return builtInError("config key not found: %s", key)
	}

	var returnValue string

	switch value.(type) {
	case string:
		returnValue = value.(string)

	case int:
		returnValue = fmt.Sprintf("%d", value.(int))

	default:
		return builtInError("keys %s has not a valid type, only string and int are allowed, got=%s", key, reflect.TypeOf(value))
	}

	return returnValue
}

func assetBuiltIn(args ...interface{}) interface{} {
	if len(args) != 1 {
		return builtInError("wrong number of arguments in asset. got=%d, want=1", len(args))
	}

	arg := args[0]

	if fmt.Sprintf("%T", arg) != "string" {
		return builtInError("argument to `asset` not supported, got %T, want=string", arg)
	}

	pathExists, path := lookForConfigKeys(govel.GetKeyFromYAML("").(map[interface{}]interface{}), "static.path")

	var pathString string

	if pathExists {
		pathString = path.(string)
	}

	s := pathString + "/" + arg.(string)

	return s
}
