package main

import "testing"

// Expected tests results have been placed here for ease of formatting.
// The tests are precise, so all spaces, tabs and line breaks must be identical.
// Note that nodes are sorted alphabetically
const (
	testRes1 = `	Errors_Test = "errors.test"
	Errors_Test2 = "errors.test2"`

	testRes2 = `	Errors_Test = "errors.test"
	Errors_Test2 = "errors.test2"

	Messages_Test = "messages.test"
	Messages_Test2 = "messages.test2"`

	testRes3 = `	Errors_Test = "errors.test"
	Errors_Test2 = "errors.test2"

	Messages_Test = "messages.test"

	Messages_SubMessages_SubMessage1 = "messages.subMessages.subMessage1"
	Messages_SubMessages_SubMessage2 = "messages.subMessages.subMessage2"`
)

func TestTransform(t *testing.T) {
	testCases := []struct {
		name            string
		expectedError   bool
		input, expected string
	}{
		{
			"blank input",
			true,
			"",
			"",
		},
		{
			"invalid json",
			true,
			"{",
			"",
		},

		{
			"empty json object",
			false,
			"{}",
			"",
		},
		{
			"single struct, single level, two members",
			false,
			`{"errors": {
				"test": "not important",
				"test2": "not important"
				}
			}`,
			testRes1,
		},
		{
			"two structs, single level, two members",
			false,
			`{"errors": {
				"test": "not important",
				"test2": "not important"
				},
			"messages": {
				"test": "not important",
				"test2": "not important"
				}
			}`,
			testRes2,
		},
		{
			"two structs, two levels, two members",
			false,
			`{"errors": {
				"test": "not important",
				"test2": "not important"
				},
			"messages": {
				"test": "not important",
				"subMessages": {
					"subMessage1": "not important",
					"subMessage2": "not important"
					}
				}
			}`,
			testRes3,
		},
	}

	for _, testCase := range testCases {
		output, err := transform([]byte(testCase.input))
		if testCase.expectedError && err == nil {
			t.Fatalf("Testing %s. Expecting error but got none",
				testCase.name,
			)
		}

		if !testCase.expectedError && err != nil {
			t.Fatalf("Testing %s. Not expecting error but got: %s",
				testCase.name,
				err,
			)
		}

		if testCase.expected != output {
			t.Errorf("Testing %s. Output not as expected. Expected: \n%s\nGot: \n%s",
				testCase.name,
				testCase.expected,
				output,
			)
		}
	}
}
