package Parser

import (
	"os"
	"regexp"
	"strings"

	px "github.com/PlayerR9/LyneCmL/Common/parsing"
)

type TokenType int

const (
	TtkTrueLit TokenType = iota
	TtkFalseLit
)

func (t TokenType) String() string {
	return [...]string{
		"true literal",
		"false literal",
	}[t]
}

var (
	numeric_literal_reg *regexp.Regexp

	string_literal_reg *regexp.Regexp

	float_literal_reg *regexp.Regexp

	command_literal_reg *regexp.Regexp

	space_reg *regexp.Regexp

	equal_sign_reg *regexp.Regexp

	long_name_reg *regexp.Regexp

	short_name_reg *regexp.Regexp

	lex_table *px.RegexTable[TokenType]
)

func init() {
	lex_table = px.NewRegexTable[TokenType]()

	// true_kw
	// 	= "1"
	// 	| "true"
	// 	| "TRUE"
	// 	| "t"
	// 	| "T"
	// 	| "True"
	// 	.
	lex_table.MustAddRegex(TtkTrueLit, `[1]|[Tt]([r][u][e])?|[T][R][U][E]`)

	// false_kw
	// 	= "0"
	// 	| "false"
	// 	| "FALSE"
	// 	| "f"
	// 	| "F"
	// 	| "False"
	// 	.
	lex_table.MustAddRegex(TtkFalseLit, `[0]|[Ff]([a][l][s][e])?|[F][A][L][S][E]`)

	// numeric_literal
	// 	= "0"
	//    | "1".."9" { "0".."9" }
	// 	.
	numeric_literal_reg = regexp.MustCompile(`^([0]|[1-9][0-9]*)`)

	// string_literal
	// 	= first_char { other_chars }
	//    | "\"" { %c - "\"" } "\""
	// 	.
	// fragment first_char
	// 	= "a".."z"
	// 	| "A".."Z"
	// 	| "0".."9"
	// 	| "_"
	// 	| "."
	// 	| "/"
	// 	| ":"
	// 	.
	// fragment other_chars
	// 	= first_char
	// 	| "-"
	//    .
	string_literal_reg = regexp.MustCompile(`^([a-zA-Z0-9_./:][a-zA-Z0-9_\-./:]*|"[^"]*")`)

	// float_literal = [ sign ] numeric_literal decimal_cmp .
	//
	// fragment decimal_cmp
	// 	= "f"
	// 	| "F"
	// 	| ( "e" | "E" ) [ sign ] numeric_literal
	// 	| "." ( "0" | { "0-9" } "1-9" )
	// 	.
	float_literal_reg = regexp.MustCompile(`^[+-]?([0]|[1-9][0-9]*)([fF]|[eE][+-]?([0]|[1-9][0-9]*)|[.]([0]|[0-9][1-9]))`)

	// command_literal = word { ( "-" | "_" ) word } .
	// fragment word = "a".."z" { "a".."z" | "0".."9" } .
	command_literal_reg = regexp.MustCompile(`^([a-z][a-z0-9]*)([\-_]word)*([a-z][a-z0-9]*)`)

	// space = " " .
	space_reg = regexp.MustCompile(`^[ ]`)

	// equal_sign = "=" .
	equal_sign_reg = regexp.MustCompile(`^[=]`)

	// long_name = "--" word { ( "-".."_" ) word } .
	// fragment word = "a".."z" { "a".."z" | "0".."9" } .
	long_name_reg = regexp.MustCompile(`^[-][-][a-z][a-z0-9]*([\-_]([a-z][a-z0-9]*))*`)

	// short_name = "-" "a".."z" { "a".."z" } .
	short_name_reg = regexp.MustCompile(`^[-][a-z]+`)
}

func Lex(data []byte) {
	// Source = command_literal Command EOF	.
	// Command
	// 	= command_literal { command_literal } [ ArgumentList ] [ FlagList ]
	// 	.
	// ArgumentList = Argument { space Argument } .
	// Argument
	// 	= bool_arg // bool
	// 	| num_lit // int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
	// 	| string_lit // string, error
	// 	| float_lit // float32, float64
	// 	.
	// bool_arg
	// 	= TRUE
	// 	| FALSE
	// 	.
	// FlagList = Flag { space Flag } .
	// Flag
	//		= LongFlag
	//		| ShortFlag
	//		.
	//
	// LongFlag = long_name [ ( equal_sign | space ) Argument ].
	// ShortFlag = short_name [ [ space ]Argument ] .
}

func Parse() {

}

func ToAst() {

}

func ParseCml() {
	line := strings.Join(os.Args, " ")

	Lex([]byte(line))

	Parse()

	ToAst()
}

const Grammar string = `
TRUE
	: [1]
	| [Tt][rue]?
	| 'TRUE'
	;

FALSE
	: [0]
	| [Ff][a-fA-F]?
	| 'FALSE'
	;

NUMBER
	: [0]
	| [1-9][0-9]*
	;

num_lit = NUMBER .

STRING : [a-zA-Z0-9_]+ | '"' .* '"' .

string_lit = STRING .

FLOAT
	: [+-]? NUMBER [.]([0] | [0-9]+[1-9]*)
	| [+-]? NUMBER [eE][+-]?[0-9]+
	| [+-]? NUMBER [fF]
	;

float_lit = FLOAT .

command = [a-z][a-z0-9]* ([_-][a-z][a-z0-9]*)* .

space = " " .
equal_sign = "=" .
dash = "-" .
long_dash = "--" .

long_flag_name = [a-z][a-z0-9]* ([_-][a-z][a-z0-9]*)* .
short_flag_name = [a-z] .

program = [a-z][a-z0-9]* ([_-][a-z][a-z0-9]*)*

Source = program Command EOF	.

Command
	= command { command } [ ArgumentList ] [ FlagList ]
	.
ArgumentList = Argument { space Argument } .
Argument
	= bool_arg // bool
	| num_lit // int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
	| string_lit // string, error
	| float_lit // float32, float64
	.

bool_arg
	= TRUE
	| FALSE
	.

FlagList = Flag { space Flag } .
Flag
	= LongFlag
	| ShortFlag
	.
LongFlag = long_dash long_flag_name ( equal_sign | space ) [ Argument ].
ShortFlag = dash short_flag_name { short_flag_name } [ space ] [ Argument ] . 
`

// any
// byte
// complex64
// complex128
// rune
// uintptr
