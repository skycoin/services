package exchange

import (
	"fmt"
	"testing"
)

func TestGetBTCValue(t *testing.T) {
	p, err := GetBTCValue()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p)
}
