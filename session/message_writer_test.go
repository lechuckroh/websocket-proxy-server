package session

import "testing"

func TestIsJSON(t *testing.T) {
	type TestData struct {
		Input    string
		IsJSON   bool
	}

	testDataList := []TestData{
		{
			Input:    "123",
			IsJSON:   false,
		},
		{
			Input:    `{"foo": "bar"}`,
			IsJSON:   true,
		},
		{
			Input:    `[1, 2]`,
			IsJSON:   true,
		},
		{
			Input:    `"[1, 2]"`,
			IsJSON:   false,
		},
	}

	for _, testData := range testDataList {
		isJSON := isJSON([]byte(testData.Input))
		if isJSON != testData.IsJSON {
			t.Errorf("isJSON mismatch. actual=%v, expected: %v", isJSON, testData.IsJSON)
		}
	}
}
