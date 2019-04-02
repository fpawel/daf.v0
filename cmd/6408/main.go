package main

import (
	"github.com/fpawel/daf/internal/assets"
)

func main(){
	assets.EnsureManifestFile()
	if err := runMainWindow(); err!=nil {
		panic(err)
	}
}
