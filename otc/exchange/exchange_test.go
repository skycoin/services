package exchange

import (
	"fmt"
	"testing"
)

func TestGetBTCPrice(t *testing.T) {
	p, err := GetBTCPrice()
	if err != nil {
		panic(err)
	}

	fmt.Println(p)
}
