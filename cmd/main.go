package main

import (
	"log"
	"text/template"

	ggen "github.com/PlayerR9/go_generator/pkg"
)

var (
	// t is the template to use.
	t *template.Template

	// Logger is the logger to use.
	Logger *log.Logger
)

func init() {
	t = template.Must(template.New("").Parse(templ))

	Logger = ggen.InitLogger("LyneCmL")
}

type GenData struct {
	PackageName string
}

func (g GenData) SetPackageName(pkg_name string) ggen.Generater {
	g.PackageName = pkg_name

	return g
}

func main() {
	err := ggen.Generate("out.go", GenData{}, t)
	if err != nil {
		panic(err)
	}
}

const templ = `
package {{ .PackageName }}

import (
	cml "github.com/PlayerR9/LyneCmL/Simple"
	cm "github.com/PlayerR9/LyneCmL"

	cmd "/cmd"
)

var (
	Program *cml.Program
)

func init() {
	Program = &cml.Program{
		Name: "my program",
		Brief: "this is the template for a program",
		Description: cm.NewDescription().
			AddLine(
				"Write here the description of the program",
			).
			Build(),
		Version: "v0.1.0",
	}

	Program.SetCommands(
		// Add here the commands of the program
	)

	err := cml.Fix(Program)
	if err != nil {
		panic(fmt.Errorf("failed to fix program: %w", err))
	}
}

func main() {
	err := cm.ExecuteProgram(Program, os.Args)
	cm.DefaultExitSequence(err)
}
`
