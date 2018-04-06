package exchange

import (
	"testing"
)

func TestGetBTCValue(t *testing.T) {
	if _, err := GetBTCValue(); err != nil {
		t.Fatal(err)
	}
}
