package main_test

import (
	"fmt"
	"io/ioutil"
	"testing"
)

const (
	sampleText1 = "Lorem ipsum dolor sit amet,\n"
	sampleText2 = "consectetur adipiscing elit,\n"
	sampleText3 = "sed do eiusmod tempor incididunt ut labore et dolore magna aliqua\n"

	expectedTestNonAppendThreeFiles = sampleText1 + sampleText2 + sampleText3
)

var (
	filenames = []string{"out-1.txt", "out-2.txt", "out-3.txt"}
)

type (
	setupData struct {
		flags     string
		filenames []string
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
		Setup      func(f *setupData)
		AssertFn   func(t *testing.T, ad assertionData)
		AssertData assertionData
	}
)

func TestBase(t *testing.T) {
	testCases := []testCase{
		{
			Name: "TestNonAppendThreeFiles",
			SetupData: &setupData{
				flags:     "--append",
				filenames: filenames,
			},
			AssertFn: verifyAssertion,
			AssertData: assertionData{
				expected: assertionItem{
					stdOutput: expectedTestNonAppendThreeFiles,
					fileOutput: map[string]string{
						filenames[0]: expectedTestNonAppendThreeFiles,
						filenames[1]: expectedTestNonAppendThreeFiles,
						filenames[2]: expectedTestNonAppendThreeFiles,
					},
				},
			},
		},
	}

	runTests(t, testCases)
}

func runTests(t *testing.T, tcs []testCase) {
	for _, test := range tcs {
		runTest(t, test)
	}
}

func runTest(t *testing.T, tc testCase) {
	t.Run(tc.Name, func(t *testing.T) {
		sd := tc.SetupData

		if tc.Setup != nil {
			tc.Setup(sd)
		}

		// result := "not calculated yet" // Execute tee command

		tc.AssertData.actual = assertionItem{
			stdOutput: "not calculated yet",
			fileOutput: map[string]string{
				filenames[0]: "not calculated yet",
				filenames[1]: "not calculated yet",
				filenames[2]: "not calculated yet",
			},
		}

		tc.AssertFn(t, tc.AssertData)
	})
}

func verifyAssertion(t *testing.T, ad assertionData) {
	t.Helper()

	if !(assertExpected(ad)) {
		t.Errorf("received value '%+v' does not match expected '%+v'\n", ad.actual, ad.expected)
	}
}

func assertExpected(ad assertionData) (ok bool) {
	// if ad.expected.stdOutput != ad.actual.stdOutput {
	// 	return false
	// }

	loadOutputFiles(&ad.expected)

	fmt.Println("----------------------")
	fmt.Printf("%+v\n", ad.actual.fileOutput)
	fmt.Println("----------------------")

	for k, expected := range ad.expected.fileOutput {
		actual, ok := ad.actual.fileOutput[k]
		if !ok || expected != actual {
			return false
		}
	}

	ok = true

	return ok
}

// Helpers
func loadOutputFiles(ai *assertionItem) (err error) {
	for _, filename := range mapKeys(ai.fileOutput) {

		content, err := ioutil.ReadFile("./test/out/" + filename)
		if err != nil {
			continue
		}

		ai.fileOutput[filename] = string(content)
	}

	return err
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
