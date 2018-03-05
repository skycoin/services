package multi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/services/coin-api/internal/locator"
	"github.com/skycoin/services/coin-api/internal/multi"
)

func TestGenerateAddress(t *testing.T) {
	loc := locator.Node{
		Host: "127.0.0.1",
		Port: 8080,
	}
	sky := multi.NewSkyService(&loc)
	rsp, err := sky.GenerateAddr(1, true)
	assert.NoError(t, err)
	if rsp.Code != 0 {
		t.Fatalf("the code must be 0, %d given", rsp.Code)
	}
}
