package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type (
	tee struct {
		input      io.ReadSeeker
		output     []io.Writer
		fileAppend bool
		fileNames  []string
		fileFlag   int
		filePerms  int
	}
)

func main() {
	append := flag.String("append", "", "append to the given FILEs, do not overwrite")
	flag.Parse()

	fmt.Println(append) // FIX: Remove this line

	ffnn := parseFileNames(flag.Args())

	t := NewTee(ffnn, false)

	t.openWiters()

	t.write()
}

func NewTee(fileNames []string, append bool) (t *tee) {
	t = &tee{
		input:      os.Stdin,
		fileAppend: append,
		fileNames:  fileNames,
		fileFlag:   os.O_CREATE | os.O_WRONLY,
		filePerms:  0644,
	}

	t.updateFileFlags()

	return t
}

func (tee *tee) updateFileFlags() {
	if tee.fileAppend {
		tee.fileFlag = tee.fileFlag | os.O_APPEND
	}

	tee.fileFlag = tee.fileFlag | os.O_TRUNC
}

func (tee *tee) openWiters() {
	for _, fn := range tee.fileNames {
		file, err := os.OpenFile(fn, tee.fileFlag, os.FileMode(tee.filePerms))
		if err != nil {
			log.Fatal(err)
		}

		tee.output = append(tee.output, file)

	}

	// Append console to output
	tee.output = append(tee.output, os.Stdout)
}

func (tee *tee) write() error {
	mw := io.MultiWriter(tee.output...)

	buffer := make([]byte, 512)

	_, err := io.CopyBuffer(mw, tee.input, buffer)
	if err != nil {
		return err
	}

	return nil
}

func parseFileNames(args []string) []string {
	// FIX: Not implemented
	fmt.Printf("%+v\n", args)
	return args
}
