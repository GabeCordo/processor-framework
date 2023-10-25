package processor

import "testing"

func TestParseInput(t *testing.T) {

	cfg := NewConfig("test")

	p, err := New(cfg)
	if err != nil {
		t.Error(err)
	}

	module, cluster, metadata, err := p.parseInput("sentiment.hourly greeting:\"hello there\" my:name\n")
	if err != nil {
		t.Error(err)
	}

	if module != "sentiment" {
		t.Error("module mismatch")
	}

	if cluster != "hourly" {
		t.Error("cluster name mismatch")
	}

	if greeting, found := metadata["greeting"]; !found && (greeting != "hello there") {
		t.Error("did not parse metadata quotations properly")
	}

	if my, found := metadata["my"]; !found && (my != "name") {
		t.Error("did not parse single word metadata (no quotations) properly")
	}
}
