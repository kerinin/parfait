package main

import (
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("parfait")
var logFormat = logging.MustStringFormatter("[parfait] %{level} %{color}%{message}%{color:reset}")

func init() {
	logging.SetFormatter(logFormat)
	logging.SetLevel(logging.DEBUG, "parfait")
}
