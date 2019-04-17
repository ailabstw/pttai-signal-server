package main

import (
	"flag"
	"net/http"

	signalserver "github.com/ailabstw/pttai-signal-server"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mattn/go-colorable"
)

func initLog(verbosity int) {
	output := colorable.NewColorableStderr()
	ostream := log.StreamHandler(output, log.TerminalFormat(true))
	glogger := log.NewGlogHandler(ostream)
	glogger.Verbosity(log.Lvl(verbosity))
	log.Root().SetHandler(glogger)
}

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")
	var verbosity = flag.Int("verbosity", 3, "verbosity")

	flag.Parse()

	initLog(*verbosity)

	server := signalserver.NewServer()

	http.HandleFunc("/signal", server.SignalHandler)

	log.Info("to ListenAndServe", "addr", *addr)
	http.ListenAndServe(*addr, nil)
}
