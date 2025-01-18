package parsing

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Schematic struct {
	Entries []*Entry `@@*`
}

type Entry struct {
	Element Element `@@`
	Command Command `| @@`
}

type Element struct {
	Type       string     `@Element`
	Properties []Property `('(' (@@ (',' @@)*)? ')')?`
	Actions    []Action   `( @@+ )?`
}

type Command struct {
	Type       string     `@Command`
	Properties []Property `('(' (@@ (',' @@)*)? ')')?`
}

type Property struct {
	Key   string `@Ident '='`
	Value string `@String`
}

type Action struct {
	Type  string  `'.' @Action`
	Units float64 `('(' @Number? ')')?`
}

var (
	SupportedElements = []string{
		"resistor", "battery", "line", "capacitor", "dot",
	}
	SupportedActions = []string{
		"right", "up", "left", "down",
	}
	SupportedCommands = []string{
		"push", "pop",
	}
)

var (
	schemGoLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Element", `(` + strings.Join(SupportedElements, "|") + `)`},
		{"Action", `(` + strings.Join(SupportedActions, "|") + `)`},
		{"Command", `(` + strings.Join(SupportedCommands, "|") + `)`},
		{"Ident", `[a-zA-Z_][a-zA-Z_0-9]*`},
		{"String", `"[^"]*"`},
		{"Number", `[-+]?[.0-9]+\b`},
		{"Punct", `\[|]|[-!()+/*=,]`},
		{"comment", `#[^\n]+`},
		{"whitespace", `\s+`},
	})
	schemGoParser = participle.MustBuild[Schematic](
		participle.Lexer(schemGoLexer),
		participle.Unquote("String"),
	)
)

func MustReadSchematic(schemFileContents []byte, schemFilePath string) (schematic *Schematic) {
	schematic, err := schemGoParser.ParseBytes(schemFilePath, schemFileContents)
	if err != nil {
		panic(err)
	}
	return schematic
}
