package ast

import (
	"bytes"
	"strings"

	"github.com/govel-framework/lamb/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }

func (vs *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")

	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) statementNode()       {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) expressionNode()      {}
func (*ExpressionStatement) statementNode()          {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) statementNode()       {}
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
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }

func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(token.LPAREN)
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(token.RPAREN)

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

func (b *Boolean) String() string {
	return b.TokenLiteral()
}

type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if(")
	out.WriteString(ie.Condition.String())
	out.WriteString(") ")

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString(token.LPAREN)
	out.WriteString(strings.Join(args, token.COMMA+" "))
	out.WriteString(token.RPAREN)

	return out.String()
}

type StringLiteral struct {
	Token  token.Token
	Value  string
	Closed bool // whether the string is closed
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
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
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
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

type MapLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *MapLiteral) expressionNode()      {}
func (hl *MapLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *MapLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}

	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type ForExpression struct {
	Token token.Token // The 'for' token
	Key   string
	Value string
	In    Expression
	Block *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fe.In.String())

	return out.String()
}

type ExtendsStatement struct {
	Token token.Token // The 'extends' token
	From  string
}

func (ee *ExtendsStatement) expressionNode()      {}
func (ee *ExtendsStatement) TokenLiteral() string { return ee.Token.Literal }
func (ee *ExtendsStatement) String() string {
	var out bytes.Buffer

	out.WriteString("extends(")
	out.WriteString(ee.From)
	out.WriteString(")")

	return out.String()
}

type SectionStatement struct {
	Token token.Token // The 'section' token
	Block *BlockStatement
	Name  string
}

func (ss *SectionStatement) expressionNode()      {}
func (ss *SectionStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SectionStatement) String() string {
	var out bytes.Buffer

	out.WriteString("section(")
	out.WriteString(ss.Name)
	out.WriteString(")")

	return out.String()
}

type DefineStatement struct {
	Token   token.Token // The 'define' token
	Name    string
	Content *BlockStatement
}

func (ds *DefineStatement) expressionNode()      {}
func (ds *DefineStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DefineStatement) String() string {
	var out bytes.Buffer

	out.WriteString("define(")
	out.WriteString(ds.Name)
	out.WriteString(")")

	return out.String()
}

type DotExpression struct {
	Token token.Token // The '.' token
	Left  Identifier
	Right Identifier
}

func (de *DotExpression) expressionNode()      {}
func (de *DotExpression) TokenLiteral() string { return de.Token.Literal }
func (de *DotExpression) String() string {
	var out bytes.Buffer

	out.WriteString(de.Left.String() + "." + de.Right.String())

	return out.String()
}

type IncludeStatement struct {
	Token token.Token // The 'include' token
	File  string
	Vars  Expression
}

func (is *IncludeStatement) expressionNode()      {}
func (is *IncludeStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncludeStatement) String() string {
	var out bytes.Buffer

	out.WriteString("include(")
	out.WriteString(is.File)

	if is.Vars != nil {
		out.WriteString(", ")
		out.WriteString(is.Vars.String())
	}

	out.WriteString(")")

	return out.String()
}

type HtmlLiteral struct {
	Token token.Token
	Value string
}

func (hl *HtmlLiteral) expressionNode()      {}
func (hl *HtmlLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HtmlLiteral) String() string       { return hl.Token.Literal }
