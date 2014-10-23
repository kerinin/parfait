package cio_lite

import (
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("cio")

func init() {
	logging.SetLevel(logging.DEBUG, "cio")
}
