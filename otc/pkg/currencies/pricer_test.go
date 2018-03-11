package currencies

import (
	"testing"
)

func TestPricerGet(t *testing.T) {
	pricer := &Pricer{
		Using: EXCHANGE,
		Sources: map[Source]*Price{
			EXCHANGE: NewPrice(100),
		},
	}

	pricer.SetSource(INTERNAL)

	if pricer.Using != INTERNAL {
		t.Fatal("set source")
	}

	_, source, _ := pricer.GetPrice()
	if source != "" {
		t.Fatal("missing source")
	}

	pricer.SetPrice(EXCHANGE, 500)
	if pricer.Sources[EXCHANGE].Amount != 500 {
		t.Fatal("set price existing")
	}

	pricer.SetPrice(INTERNAL, 20)
	if pricer.Sources[INTERNAL].Amount != 20 {
		t.Fatal("set price new")
	}
}
