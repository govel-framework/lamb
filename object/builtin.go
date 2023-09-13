package object

type BuiltinFunction func(args ...interface{}) interface{}

type Builtin struct {
	Fn BuiltinFunction
}
