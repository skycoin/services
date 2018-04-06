package otc

import (
	"testing"
)

func TestConfigNew(t *testing.T) {
	conf, err := NewConfig("config_test.toml")

	if err != nil {
		t.Fatal(err)
	}

	if conf.SKY.Node != "test" {
		t.Fatalf(`expected "test", got "%s"`, conf.SKY.Node)
	}
}

func TestConfigNewBad(t *testing.T) {
	if _, err := NewConfig("bad"); err == nil {
		t.Fatal("should return error")
	}
}
