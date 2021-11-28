package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
	append := flag.Bool("append", false, "append to the given FILEs, do not overwrite")
	flag.Parse()

	ffnn := filenames(flag.Args())

	t := NewTee(ffnn, *append)

	t.openWriters()

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

func (tee *tee) openWriters() {
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

func filenames(args []string) []string {
	ffnn := []string{}

	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			ffnn = append(ffnn, a)
		}
	}

	fmt.Printf("%+v\n", ffnn)
	return ffnn
}
