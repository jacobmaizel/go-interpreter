package parser

import (
	"fmt"

	"github.com/jacobmaizel/go-interpreter/ast"
	"github.com/jacobmaizel/go-interpreter/lexer"
	"github.com/jacobmaizel/go-interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

// uses a lexer to initialize a new parser, and sets the initial curToken and peekTokens
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// read two tokens to set curToken & peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// returns the stored errors in the parser
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got=%s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// advances the parsers token values
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// parses a list of statements from the parser
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// creates an ast node for this return statement, advances the parsers's tokens
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: skipping expressions, until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// ie. let varName = 5;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// expects a name for the variable
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Creates the identifer node for this statement
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// ensures we have an = sign
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// here we store the expression ie. 5 from the func comment above. could be
	// anything resulting in a value like 3*5 or a+b
	// TODO: skipping expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// check if the token the parser is currently pointing to matches the give token
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// check if the given token matches the type of the peekToken (next) on the parser
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// if the given token is what we expect it to be, advance the parser to the next
// token and return true, else false
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
