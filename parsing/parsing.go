package parsing

import (
	"os"

	"github.com/alecthomas/participle/v2"
)

type Schematic struct {
	Components []*Component `@@*`
}

type Component struct {
	Name       string     `@Ident "("`
	Properties []Property `(@@ ("," @@)*)? ")"`
}

type Property struct {
	Key   string `@Ident "="`
	Value Value  `@@`
}

type Value interface{ value() }

type String struct {
	String string `@String`
}

func (String) value() {}

type Number struct {
	Number float64 `@Float`
}

func (Number) value() {}

func ReadSchematic(schematicFilePath string) (schematic *Schematic, err error) {
	parser, err := participle.Build[Schematic](
		participle.Unquote("String"),
		participle.Union[Value](String{}, Number{}),
	)

	if err != nil {
		panic(err)
	}

	schemFile, err := os.Open(schematicFilePath)
	if err != nil {
		exit("Error opening schemFile: %s\n", err)
	}

	return parser.Parse(schematicFilePath, schemFile)
}
