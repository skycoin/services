package exchange

import (
	"testing"
)

func TestGetBTCPrice(t *testing.T) {
	_, err := GetBTCValue()
	if err != nil {
		t.Fatal(err)
	}
}
