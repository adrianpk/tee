package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	testOutPath = "./test/out/"
)

const (
	sampleTextIn = `Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
`
)

var (
	files = []string{"out-1.txt", "out-2.txt", "out-3.txt"}
)

type (
	setupData struct {
		append     bool
		files      []string
		stdin      *os.File
		stdout     *os.File
		stdoutChan chan string
		prevStdin  *os.File
		prevStdout *os.File
	}

	assertionData struct {
		actual   assertionItem
		expected assertionItem
	}

	assertionItem struct {
		stdOutput  string
		fileOutput map[string]string
	}

	testCase struct {
		Name       string
		SetupData  *setupData
		Setup      func(t *testing.T, f *setupData) error
		AssertFn   func(t *testing.T, ad *assertionData)
		AssertData *assertionData
	}
)

// NOTE: The intention is to create a table driven test but at the moment
//  only one case is being tested (happy path).
// Add more cases.
func TestBase(t *testing.T) {
	updateOutputFilenames()

	testCases := []testCase{
		{
			Name: "TestNonAppendThreeFiles",
			SetupData: &setupData{
				append: true,
				files:  files,
			},
			Setup:    runCommand,
			AssertFn: verifyAssertion,
			AssertData: &assertionData{
				expected: assertionItem{
					stdOutput: sampleTextIn,
					fileOutput: map[string]string{
						files[0]: sampleTextIn,
						files[1]: sampleTextIn,
						files[2]: sampleTextIn,
					},
				},
			},
		},
	}

	runTests(t, testCases)
}

func runTests(t *testing.T, tcs []testCase) {
	for _, test := range tcs {
		setup(t)

		runTest(t, test)
	}
}

func runTest(t *testing.T, tc testCase) {
	t.Run(tc.Name, func(t *testing.T) {
		sd := tc.SetupData

		if tc.Setup != nil {
			err := tc.Setup(t, sd)
			if err != nil {
				t.Fatalf("test setup error: %s", err.Error())
			}
		}

		// Init results
		tc.AssertData.resetActual()

		// Read files
		loadOutputFiles(&tc)

		// Restore original in/out
		tc.SetupData.restoreStdin(t)
		tc.SetupData.restoreStdout(t)

		// Read stdout
		readStdout(&tc)

		// Assert
		tc.AssertFn(t, tc.AssertData)
	})
}

func runCommand(t *testing.T, sd *setupData) (err error) {
	// In
	stdin, prev := mockStdin(t)
	sd.stdin = stdin
	sd.prevStdin = prev

	// Out
	stdout, prev, outChan := mockStdout(t)
	sd.stdout = stdout
	sd.prevStdout = prev
	sd.stdoutChan = outChan

	tee := NewTee(sd.files, sd.append)
	err = tee.execute()

	mockUserInput()

	return err
}

func verifyAssertion(t *testing.T, ad *assertionData) {
	t.Helper()

	if !(assertExpected(ad)) {
		t.Errorf("received value '%+v' does not match expected '%+v'\n", ad.actual, ad.expected)
	}
}

func assertExpected(ad *assertionData) (ok bool) {
	if ad.expected.stdOutput != ad.actual.stdOutput {
		return false
	}

	for k, expected := range ad.expected.fileOutput {
		actual, ok := ad.actual.fileOutput[k]
		if !ok || expected != actual {
			return false
		}
	}

	ok = true

	return ok
}

func (ad *assertionData) resetActual() {
	ad.actual.stdOutput = ""
	ad.actual.fileOutput = map[string]string{}
}

func (sd *setupData) restoreStdin(t *testing.T) {
	os.Stdin = sd.prevStdin

	err := sd.stdin.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func (sd *setupData) restoreStdout(t *testing.T) {
	os.Stdout = sd.prevStdout

	err := sd.stdout.Close()
	if err != nil {
		t.Errorf(err.Error())
	}

	// close(sd.stdoutChan)
}

// Setup & teardown
func setup(t *testing.T) {
	err := os.RemoveAll(testOutPath)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = os.Mkdir(testOutPath, 0777)
	if err != nil {
		t.Errorf(err.Error())
	}
}

// Helpers
func updateOutputFilenames() {
	for i, _ := range files {
		files[i] = testOutPath + files[i]
	}
}

func readStdout(tc *testCase) {
	ad := tc.AssertData
	ad.actual.stdOutput = <-(tc.SetupData.stdoutChan)
}

func loadOutputFiles(tc *testCase) (err error) {
	ad := tc.AssertData
	for _, filename := range mapKeys(ad.expected.fileOutput) {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			continue
		}

		ad.actual.fileOutput[filename] = string(content)
	}

	return err
}

func mockUserInput() {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			break
		}

		fields := strings.Fields(line)
		fmt.Println(fields)
	}

	if err := s.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func mockStdin(t *testing.T) (tmp, old *os.File) {
	textBytes := []byte(sampleTextIn)

	tmpStdin, err := ioutil.TempFile("", "in")
	if err != nil {
		log.Fatal(err)
	}

	_, err = tmpStdin.Write(textBytes)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tmpStdin.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	old = os.Stdin
	os.Stdin = tmpStdin

	return tmpStdin, old
}

func mockStdout(t *testing.T) (tmp, old *os.File, tmpOutChan chan string) {
	old = os.Stdout
	read, tmp, _ := os.Pipe()
	os.Stdout = tmp

	print()

	outChan := make(chan string)

	go func() {
		var b bytes.Buffer
		io.Copy(&b, read)
		outChan <- b.String()
	}()

	return tmp, old, outChan
}

func mapKeys(m map[string]string) (keys []string) {
	keys = make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}
