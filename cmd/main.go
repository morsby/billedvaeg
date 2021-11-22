package main

import (
	"time"

	"github.com/morsby/billedvaeg"
)

func main() {
	billedvaeg.Server(5000, 15*time.Second)
}
