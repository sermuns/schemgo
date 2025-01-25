package parsing

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/sermuns/schemgo/drawing"
	"maps"
	"slices"
	"strings"
)

type Schematic struct {
	Entries []Entry `@@*`
}

type Entry struct {
	Element Element `@@`
	Command Command `| @@`
}

type Element struct {
	Type       string     `@Element`
	Properties []Property `('(' (@@ (',' @@)*)? ')')?`
	Movements  []Movement `( @@+ )?`
	Labels     []Label    `( ".label" '(' @@ ')' )*`
}

type Command struct {
	Type       string     `@Command`
	Properties []Property `('(' (@@ (',' @@)*)? ')')?`
}

type Property struct {
	Key   string `@Ident '='`
	Value string `@String`
}

type Movement struct {
	Type  string  `'.' @Movement`
	Units float64 `( '(' @Number ')')?`
}

type Label struct {
	Content string `@String`
	// Placement string `(, @Ident)?`
}

func mapKeysStringJoin[V any](m map[string]V) string {
	return strings.Join(slices.Collect(maps.Keys(m)), "|")
}

var (
	schemGoLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Element", `(` + mapKeysStringJoin(drawing.ElemTypeToRenderFunc) + `)`},
		{"Command", `(` + mapKeysStringJoin(drawing.CommandTypeToFunc) + `)`},
		{"Movement", `(` + mapKeysStringJoin(drawing.MovementFuncs) + `)`},
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
