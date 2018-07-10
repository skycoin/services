package active

import (
	"testing"
)

func TestToken(t *testing.T) {
	instance := newDockerHub("skycoin/skycoin")
	instance.getToken()
	//t.Log(fmt.Sprintf("%+v", instance.token))
	if instance.token == nil {
		t.Fail()
	}
}
