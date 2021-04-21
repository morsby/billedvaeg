package main

import (
	"html/template"
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/morsby/billedvaeg"
	"github.com/morsby/billedvaeg/web"
)

func main() {
	web.Serve(8000)
}

func init() {
	genIndexPage()
	genScripts()
}

type pageData struct {
	Positions billedvaeg.Positions
}

func genIndexPage() {
	tpl := template.Must(template.ParseFiles("web/index.gohtml"))

	err := os.MkdirAll("web/dist", 777)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("web/dist/index.html")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	data := pageData{
		Positions: billedvaeg.Positions{}.FromJSON(),
	}

	tpl.Execute(f, data)
	return
}

func genScripts() {
	result := api.Build(api.BuildOptions{
		EntryPoints:      []string{"web/scripts.ts"},
		Bundle:           true,
		Format:           api.FormatIIFE,
		MinifyWhitespace: true,
		MinifySyntax:     true,
		Write:            true,
		Outdir:           "web/dist",
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}
