package ast

import (
	"bytes"

	"github.com/jacobmaizel/go-interpreter/token"
)

type Node interface {
	TokenLiteral() string // mainly debugging and testing
	String() string
}
type Statement interface {
	Node
	statementNode() // helps us to guide compiler and knowing when we misuse statement in place of expression and vice versa
}
type Expression interface {
	Node
	expressionNode()
}

// Program will be the root node of every AST our parser produces
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

// ie. let asdf = 5 * 6;
type LetStatement struct {
	Token token.Token // token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
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

// ie. the a in var a=5;
type Identifier struct {
	Token token.Token // token.IDENT token
	Value string
}

func (i *Identifier) String() string { return i.Value }

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // 'return token'
	ReturnValue Expression
}

// ie. return a+b, return 5, return abc
func (rs *ReturnStatement) statementNode()       {}
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

// ie. x+5;
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

// ie. return a+b, return 5, return abc
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
