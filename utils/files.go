package utils

import (
	"log"
	"os"
	"path/filepath"
)

func SaveFile(name string, contents []byte) error {
	file, err := os.Create(name)
	if err != nil {
		log.Println("Error creating file", err)

		return err
	}

	defer file.Close()

	if _, err := file.Write(contents); err != nil {
		log.Println("Error writting file", name, "; ", err)

		return err
	}

	return nil
}

func Dump(name, contents string) {
	fpath := filepath.Join("dumps/", name)

	f, err := os.Create(fpath)
	if err != nil {
		log.Printf("Failed to create file %s\n", fpath)

		return
	}

	defer f.Close()

	if written, err := f.WriteString(contents); err != nil {
		log.Printf("Error writing to file %v %s\n", written, fpath)

		return
	}
}
