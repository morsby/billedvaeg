package web

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/morsby/billedvaeg"
)

func Compile(w io.Writer) error {
	script := genScripts()
	return genIndexPage(w, string(script))
}

type pageData struct {
	Positions billedvaeg.Positions
	Script    template.JS
}

//go:embed index.gohtml
var tmpl string

func genIndexPage(w io.Writer, script string) error {
	tpl := template.Must(template.New("index").Parse(tmpl))
	data := pageData{
		Positions: billedvaeg.Positions{}.FromJSON(),
		Script:    template.JS(script),
	}

	return tpl.Execute(w, data)
}

//go:embed scripts.ts
var script string

func genScripts() []byte {
	script = strings.Replace(script, `import positions from "../positions.json";`, fmt.Sprintf("const positions = %s", billedvaeg.PositionsJson), -1)

	result := api.Build(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents: script,
			Loader:   api.LoaderTS,
		},
		Bundle:            true,
		Format:            api.FormatIIFE,
		MinifyWhitespace:  true,
		MinifySyntax:      true,
		MinifyIdentifiers: true,
		Write:             false,
		Outdir:            "/",
	})

	if len(result.Errors) > 0 {
		fmt.Println(result.Errors)
		os.Exit(1)
	}

	return result.OutputFiles[0].Contents
}
