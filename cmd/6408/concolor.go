package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-logfmt/logfmt"
	"path/filepath"
)

type colorStdOutWriter struct{}

func (x colorStdOutWriter) Write(p []byte) (int, error) {
	d := logfmt.NewDecoder(bytes.NewReader(p))

	var (
		msg, time, fields, fileLine string
	)
	lev := "???"
	c := color.New(color.FgWhite)
	for d.ScanRecord() {
		for d.ScanKeyval() {
			value := string(d.Value())
			key := string(d.Key())
			switch key {
			case "msg":
				msg = value
			case "file":
				fileLine = filepath.Base(value)
			case "time":
				time = value
			case "level":
				switch value {
				case "error", "panic", "fatal":
					c = color.New(color.FgHiRed, color.Bold)
					lev = "ERR"
				case "warn", "warning":
					c = color.New(color.FgHiMagenta, color.Bold)
					lev = "WRN"
				case "info":
					c = color.New(color.FgHiWhite, color.Bold)
					lev = "INF"
				case "debug":
					c = color.New(color.FgWhite)
					lev = "DBG"

				}
			case "func":

			default:
				fields += fmt.Sprintf(" %s=%q", key, value)
			}
		}
	}
	if len(fields) > 0 {
		msg += " " + fields
	}
	_, _ = color.New(color.FgGreen).Fprint(color.Output, time+" ")
	_, _ = c.Fprint(color.Output, lev+" ")
	_, _ = color.New(color.FgCyan).Fprint(color.Output, fileLine+"\t")

	return c.Fprintln(color.Output, msg)
}
