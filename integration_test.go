package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"gh-comment": main1,
	}))
}

func TestIntegration(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
		Setup: func(env *testscript.Env) error {
			// Set up test environment
			env.Setenv("GH_TOKEN", "test-token")
			env.Setenv("GH_HOST", "github.com")
			return nil
		},
		Condition: func(cond string) (bool, error) {
			switch cond {
			case "has-gh":
				// Check if gh CLI is available
				_, err := os.Stat("/usr/local/bin/gh")
				return err == nil, nil
			default:
				return false, nil
			}
		},
	})
}

// main1 is a wrapper around main that returns an exit code
func main1() int {
	// Call the actual main function and handle panics
	defer func() {
		if r := recover(); r != nil {
			// If main panics, return error code
			os.Exit(1)
		}
	}()
	
	// Call the real main function
	main()
	return 0
}
