package main

import (
	"github.com/fpawel/daf/internal/data"
	"github.com/powerman/structlog"
	"os"
	"path/filepath"
	_ "runtime/cgo"
)

func main() {

	structlog.DefaultLogger.
		SetPrefixKeys(
			structlog.KeyApp, structlog.KeyPID, structlog.KeyLevel, structlog.KeyTime,
		).
		SetDefaultKeyvals(structlog.KeyApp, filepath.Base(os.Args[0])).
		SetSuffixKeys(
			structlog.KeyStack,
		).
		SetKeysFormat(map[string]string{
			structlog.KeyTime:   " %[2]s",
			structlog.KeySource: " %6[2]s",
			KeyAddr:             " %[2]X",
			KeyEN6408:           " %+v",
		}).SetTimeFormat("15:04:05")

	log := structlog.New()
	log.Info("start", structlog.KeyTime, now())

	data.Open()

	if err := runMainWindow(); err != nil {
		panic(err)
	}

}
