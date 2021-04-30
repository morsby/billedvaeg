package main

import (
	"os"

	"github.com/morsby/billedvaeg/web"
)

func main() {
	web.Compile(os.Stdout)
}
