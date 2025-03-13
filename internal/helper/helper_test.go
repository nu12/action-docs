package helper

import (
	"testing"
)

func TestHash(t *testing.T) {
	var data = "data"
	result := Hash(data)
	expected := "8d777f385d3dfec8815d20f7496026dc"
	if expected != result {
		t.Errorf("error: %s", "Hash doesn't match")
	}
}
