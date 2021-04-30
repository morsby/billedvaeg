package web

import (
	"html/template"
	"io"
	"os"

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

func genIndexPage(w io.Writer, script string) error {
	tpl := template.Must(template.ParseFiles("web/index.gohtml"))
	data := pageData{
		Positions: billedvaeg.Positions{}.FromJSON(),
		Script:    template.JS(script),
	}

	return tpl.Execute(w, data)
}

func genScripts() []byte {
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{"web/scripts.ts"},
		Bundle:            true,
		Format:            api.FormatIIFE,
		MinifyWhitespace:  true,
		MinifySyntax:      true,
		MinifyIdentifiers: true,
		Write:             false,
		Outdir:            "/",
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	return result.OutputFiles[0].Contents
}
