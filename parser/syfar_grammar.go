package parser

import "github.com/alecthomas/participle/v2/lexer"

// Package parser implements a parser for Syfar Test Language syntax.

type SyfarFile struct {
	Entries []*Entry `parser:"@@*"`
}

type Entry struct {
	Stepper       *Stepper       `parser:"@@"`
	Action        *Action        `parser:"|@@"`
	TestSet       *TestSet       `parser:"|@@"`
	Test          *Test          `parser:"|@@"`
	Variable      *Variable      `parser:"|@@"`
	MultiVariable *MultiVariable `parser:"|@@"`
	VarSet        *VarSet        `parser:"|@@"`
	SecretSet     *SecretSet     `parser:"|@@"`
	Print         *Print         `parser:"|@@"`
	Import        *Import        `parser:"|@@"`
}

type Action struct {
	Type       string `parser:"'action' @(Ident|String)"`
	Id         string `parser:" @(Ident|String)"`
	Prefix     *string
	Attributes []*ActionAttribute `parser:" '{' @@* '}'"`
	Pos        lexer.Position
}

type ActionAttribute struct {
	Parameter *Assignment `parser:"@@"`
	Test      *Test       `parser:"|@@"`
	TestSet   *TestSet    `parser:"|@@"`
	Out       *Out        `parser:"|@@"`
}

type Stepper struct {
	Id    string   `parser:"'steps' @(Ident|String)"`
	Steps []*Steps `parser:" '{' @@* '}'"`
}

type Steps struct {
	Action *Action `parser:"@@"`
}

type TestSet struct {
	Description string  `parser:"'tests' @String"`
	Tests       []*Test `parser:" '{' @@* '}'"`
}

type Test struct {
	Description  string         `parser:"'test' @String"`
	Expectations []*Expectation `parser:" '{' @@* '}'"`
}

type Expectation struct {
	Items []*ExpectationItem `parser:"'expect' '{' @@* '}'"`
}

type ExpectationItem struct {
	Symbolic *SymbolicCheck `parser:"@@"`
	Chain    *ChainCheck    `parser:"|@@"`
}

type SymbolicCheck struct {
	Key   string `parser:"@Ident (@'[' @Int @']')? ( @'.' @Ident (@'[' @Int @']')? )*"`
	Opp   string `parser:"@('=='|'<='|'>='|'<'|'>'|'!='|'eq'|'gt'|'lt'|'le'|'ge'|'ne')"`
	Value *Value `parser:"@@"`
}

type ChainCheck struct {
	Key   string       `parser:"@Ident ( @'.' @Ident )* ':' "`
	Chain []*ChainItem `parser:"@@*"`
}

type ChainItem struct {
	Start  string   `parser:"@'.'? @'to'"`
	Negate bool     `parser:"(@'.' @'not')?"`
	Deep   bool     `parser:"(@'.' @'deep')?"`
	Method string   `parser:" @'.' @('be')  @'.' @('eq'|'gt'|'lt'|'le'|'ge'|'ne')"`
	Args   []*Value `parser:"('(' @@* ')')?"`
}

type ChainMethod struct {
	Name string `parser:"@('be' |'eq'|'gt'|'lt'|'le'|'ge'|'ne'| 'been' | 'is' | 'that' | 'which' | 'and' | 'has' | 'have' | 'with' | 'at' | 'of' | 'same' | 'but' | 'does' | 'still' | 'also' | 'not' | 'deep' )"`
}

type ChainArg struct {
	Value Value `parser:"@Ident"`
}

type Variable struct {
	Name  string `parser:"'var' @Ident"`
	Type  string `parser:"(':' @('number'|'string'|'array'|'bool'|'object'))?"`
	Value *Value `parser:"'=' @@"`
	Pos   lexer.Position
}

type MultiVariable struct {
	Variables []*Assignment `parser:"'var' '(' @@* ')'"`
}

type VarSet struct {
	Id        string        `parser:"'vars' @(Ident|String)"`
	Variables []*Assignment `parser:" '{' @@* '}'"`
	Pos       lexer.Position
}

type Assignment struct {
	Name  string `parser:"@Ident"`
	Value *Value `parser:"'=' @@"`
	Pos   lexer.Position
}

type SecretSet struct {
	Id        string        `parser:"'secrets' @(Ident|String)"`
	Variables []*Assignment `parser:" '{' @@* '}'"`
	Pos       lexer.Position
}

type Bool bool

func (b *Bool) Capture(v []string) error { *b = v[0] == "true"; return nil }

type Value struct {
	Boolean    *Bool    `parser:" @('true'|'false')"`
	Identifier *string  `parser:"| @Ident (@'[' @Int @']')? ( @'.' @Ident (@'[' @Int @']')? )*"`
	String     *string  `parser:"| @(String|Char|RawString)"`
	Number     *float64 `parser:"| @(Float|Int)"`
	Array      []*Value `parser:"| '[' ( @@ ','? )* ']'"`
	Json       *JSON    `parser:"|@@"`
	Map        map[string]interface{}
	Any        interface{}
}

type JSON struct {
	Attributes []*JSONAttribute `parser:"'{' @@* '}'"`
}

type JSONAttribute struct {
	Name  string `parser:"@(Ident|String)"`
	Value *Value `parser:"':' @@ (',')?"`
	Pos   lexer.Position
}

type HeaderValue struct {
	Name  string `parser:"@(Ident|String)"`
	Value *Value `parser:"':' @@ (',')?"`
}

type VariableType struct {
	Name string `parser:"'number' | 'string'"`
}

type Print struct {
	Id        string   `parser:"'print'"`
	Variables []*Value `parser:"'{' (@@ ','?)* '}'"`
	Pos       lexer.Position
}

type Out struct {
	Variables []*OutAssignment `parser:"'out' '{' (@@ ','?)* '}'"`
}

type OutAssignment struct {
	Name       string `parser:"@Ident"`
	Identifier string `parser:"'=' @Ident (@'[' @Int @']')? ( @'.' @Ident (@'[' @Int @']')? )*"`
	Pos        lexer.Position
}

type Import struct {
	Files []string `parser:"'import' '(' (@String ','?)* ')'"`
}

type Argument struct {
	Name    string `parser:"@Ident"`
	Type    string `parser:"':' @('Id'|'Value')"`
	Default *Value `parser:"( '=' @@ )?"`
}
