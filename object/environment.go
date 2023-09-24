package object

import "github.com/govel-framework/lamb/token"

func NewEnvironment() *Environment {
	s := make(map[string]interface{})
	return &Environment{store: s, outer: nil, ExtendsFrom: parentTemplate{
		Sections: make(map[string]SectionContent),
	}}
}

func CopyEnvironment(env *Environment) *Environment {
	newEnv := NewEnvironment()
	newEnv.store = env.store
	newEnv.outer = env.outer
	newEnv.ExtendsFrom = env.ExtendsFrom

	s, _ := env.Get("sessions")

	newEnv.Set("sessions", s)

	return newEnv
}

type SectionContent struct {
	Token   token.Token // The token of the section.
	Name    string      // The name of the section.
	Content interface{} // The default or real content of the section.
}

type parentTemplate struct {
	Sections map[string]SectionContent // The sections in the template.
	From     string                    // The template that extends from.
}

type Environment struct {
	store    map[string]interface{}
	outer    *Environment
	FileName string

	InExtends bool
	IsExtends bool
	InSection bool
	InDefine  bool

	ExtendsFrom parentTemplate // The template that extends from.
}

func (e *Environment) Get(name string) (interface{}, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}
func (e *Environment) Set(name string, val interface{}) {
	e.store[name] = val
}

func (e *Environment) Delete(name string) {
	delete(e.store, name)
}
