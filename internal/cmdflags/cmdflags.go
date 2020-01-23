package cmdflags

import (
	"flag"
)

var DevMode = flag.Bool("devmode", false, "Start server in development mode, enabling debug interface and pretty logs")

func Parse() {
	flag.Parse()
}
