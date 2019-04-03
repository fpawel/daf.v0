package main

import (
	"github.com/fpawel/daf/internal/assets"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "15:04:05"})
	logrus.SetOutput(colorStdOutWriter{})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)

	assets.EnsureManifestFile()
	if err := runMainWindow(); err != nil {
		panic(err)
	}
}
