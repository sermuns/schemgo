package parsing

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
)

type Schematic struct {
	Elements []*Element `@@*`
}

type BoundingBox struct {
	xmin, xmax int
	ymin, ymax int
}

type Element struct {
	Type        string     `@( "battery" | "line" | "resistor" )`
	Properties  []Property `('(' (@@ (',' @@)*)? ')')?`
	Actions     []Action   `( @@+ )?`
	BoundingBox BoundingBox
}

type Property struct {
	Key   string `@Ident "="`
	Value Value  `@@`
}

type Action struct {
	Type  string `'.' @Ident`
	Value int    `('(' @Int? ')')?`
}

type Value interface{ value() }

type String struct {
	String string `@String`
}

func (String) value() {}

type Number struct {
	Number int `@Int`
}

func (Number) value() {}

func ReadSchematic(schematicFilePath string) (schematic *Schematic, err error) {
	parser, err := participle.Build[Schematic](
		participle.Unquote("String"),
		participle.Union[Value](String{}, Number{}),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed parsing: %s\n", err)
		os.Exit(1)
	}

	schemFile, err := os.Open(schematicFilePath)
	if err != nil {
		panic(err)
	}

	return parser.Parse(schematicFilePath, schemFile)
}
