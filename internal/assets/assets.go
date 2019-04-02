package assets

import (
	"log"
	"os"
	"path/filepath"
)

func EnsureManifestFile(){

	fileName := filepath.Base(os.Args[0]) + ".manifest"
	fileName = filepath.Join(filepath.Dir(os.Args[0]), fileName)
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := Asset("assets/exe.manifest")
	if err != nil{
		log.Fatal(err)
	}
	if _, err := file.Write(bytes); err != nil {
		log.Fatal(err)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}
