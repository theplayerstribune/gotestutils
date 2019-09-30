package testutils

import (
	"io/ioutil"
	"os"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

// LoadYAMLEnv loads a "dotenv" file structured as YAML into the environment.
// It preserves the original environment, which can be reset by calling the
// "resetEnv" function it returns. Eg:
//
// 		func TestExample(t *testing.T) {
// 			resetEnv := gotestutils.LoadYAMLEnv(t)
// 			defer resetEnv()
//
// 			// Run tests here....
// 		}
//
func LoadYAMLEnv(t *testing.T, filepath string) (resetEnv func()) {
	t.Helper()

	if filepath == "" {
		filepath = "./env.yaml"
	}

	byt, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	oldenvvars := map[string]string{}
	newenvvars := map[string]string{}
	yaml.Unmarshal(byt, &newenvvars)

	for name, value := range newenvvars {
		oldenvvars[name] = os.Getenv(name)
		os.Setenv(name, value)
	}

	return func() {
		for name := range newenvvars {
			os.Setenv(name, oldenvvars[name])
		}
	}
}
