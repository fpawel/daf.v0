package main

import (
	"github.com/powerman/structlog"
	"time"
)

const (
	KeyAddr   = "_addr"
	KeyPlace  = "_place"
	KeyEN6408 = "_en6408"
)

var log = structlog.DefaultLogger

func now() string {
	return time.Now().Format("15:04:05")
}
