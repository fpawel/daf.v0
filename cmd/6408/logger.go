package main

import (
	"github.com/fatih/color"
	"log"
	"os"
)

func printlnWithTimeGreen(attrs ...color.Attribute) func(args ...interface{}) {
	return func(args ...interface{}) {
		printTimeGreen()
		color.Set(attrs...)
		defer color.Unset()
		aLogger.Println(args...)
	}
}

func printWithTimeGreen(attrs ...color.Attribute) func(args ...interface{}) {
	return func(args ...interface{}) {
		printTimeGreen()
		color.Set(attrs...)
		defer color.Unset()
		aLogger.Print(args...)
	}
}

func printfWithTimeGreen(attrs ...color.Attribute) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		printTimeGreen()
		color.Set(attrs...)
		defer color.Unset()
		aLogger.Printf(format, args...)
	}
}

func printTimeGreen() {
	color.Set(color.FgHiGreen)
	loggerTime.Print(" ")
}

var (
	loggerTime = log.New(os.Stdout, "", log.Ltime)
	aLogger    = log.New(os.Stdout, "", 0)

	info   = printWithTimeGreen(color.FgHiCyan)
	infof  = printfWithTimeGreen(color.FgHiCyan)
	infoln = printlnWithTimeGreen(color.FgHiCyan)

	perror  = printWithTimeGreen(color.FgRed, color.Bold)
	errorf  = printfWithTimeGreen(color.FgRed, color.Bold)
	errorln = printlnWithTimeGreen(color.FgRed, color.Bold)
)
