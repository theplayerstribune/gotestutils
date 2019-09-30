package testutils

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/golden"
)

const (
	// GoldenTestFunctionSuffix is the shared test name suffix for all
	// tests that run against golden files, for ease in filtering them.
	GoldenTestFunctionSuffix = "Snapshot"

	// GoldenDir is the (relative to the test file) directory name for
	// where we store golden files.
	GoldenDir = "testdata"

	// GoldenUpdateFlag is the flag that must be passed to `go test` to
	// update the golden files.
	GoldenUpdateFlag = "test.update-golden"

	// JSONIndent is the number of spaces to indent JSON in the golden
	// files.
	JSONIndent = 2
)

var (
	// JSONEncoder stores specific rules for serializing JSON.
	JSONEncoder *json.Encoder

	// JSONFileBuffer is a buffer that by is attached to the JSONEncoder.
	JSONFileBuffer bytes.Buffer
)

func init() {
	JSONEncoder = json.NewEncoder(&JSONFileBuffer)
	JSONEncoder.SetIndent("", indent(JSONIndent))
}

// AssertMatchesGolden is a simple wrapper around gotest.tools/golden that:
//
//  * prepends "golden." to the golden file's name,
//  * ensures standard naming of golden tests,
//  * ensures standard placement of golden tests,
//
func AssertMatchesGolden(
	t *testing.T,
	actualValueAsGoldenFile,
	expectedValue string,
) {
	t.Helper()

	if _, err := os.Stat(GoldenDir); os.IsNotExist(err) {
		if hasFlag(GoldenUpdateFlag) {
			os.Mkdir(GoldenDir, os.ModePerm)
		}
	}

	if !strings.HasSuffix(t.Name(), GoldenTestFunctionSuffix) {
		panic(fmt.Sprintf(
			"snapshot test is named %q: must end %q",
			t.Name(),
			"..."+GoldenTestFunctionSuffix,
		))
	}

	golden.Assert(t, expectedValue, "golden."+actualValueAsGoldenFile)
}

// AssertMarshaledJSONGolden is a simple wrapper around AssertMatchesGolden
// that:
//
//  * (de)serializes the values as JSON,
//  * uses a JSON encoder with a standard indentation to make the files
//    and diffs easier to read.
//
func AssertMarshaledJSONGolden(
	t *testing.T,
	actualValueAsGoldenFile string,
	expectedValue interface{},
) {
	t.Helper()

	err := JSONEncoder.Encode(expectedValue)
	require.NoError(t, err)

	str := JSONFileBuffer.String()
	JSONFileBuffer.Reset()

	AssertMatchesGolden(t, actualValueAsGoldenFile, str)
}

func indent(n int) string {
	indentation := make([]rune, n)
	for i := 0; i < n; i++ {
		indentation[i] = ' '
	}
	return string(indentation)
}

func hasFlag(name string) (found bool) {
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return
}
