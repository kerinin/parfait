package main

import (
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("parfait")
var logFormat = logging.MustStringFormatter("%{level} %{color}%{message}%{color:reset} [%{shortfile}]")

func init() {
	logging.SetFormatter(logFormat)
	logging.SetLevel(logging.DEBUG, "parfait")
}
