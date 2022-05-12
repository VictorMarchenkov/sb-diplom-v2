package main

import (
	"flag"
	"sb-diplom-v2/internal/app"
)

var port = flag.Int("p", 8282, "server port")

// main -.
func main() {
	flag.Parse()
	app.Run(*port)
}
