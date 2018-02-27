package exchange

import (
	"fmt"
	"testing"
)

func TestGetBTCPrice(t *testing.T) {
	p, err := GetBTCValue()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p)
}
