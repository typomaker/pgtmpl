package pgtemplate_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/cryomator/pgtemplate"
)

func printQuery(q pgtemplate.Query) {
	fmt.Print("TEXT: ")
	fmt.Println(q.String())
	fmt.Print("ARGUMENT: ")
	fmt.Println(q.Args()...)
}

// temporaryFile defines the contents of a template to be stored in a file, for testing.
type temporaryFile struct {
	name     string
	contents string
}

func makeTemporaryFile(file temporaryFile) string {
	dir, err := os.MkdirTemp("", "template")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(filepath.Join(dir, file.name))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = io.WriteString(f, file.contents)
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
func makeTemporaryDirectory(files []temporaryFile) string {
	dir, err := os.MkdirTemp("", "template")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		f, err := os.Create(filepath.Join(dir, file.name))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = io.WriteString(f, file.contents)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}
