package gen

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestAttrs(t *testing.T) {
	b := `
attrs:
  a: 1
  b: "2"
`
	cfg := &Config{}

	err := yaml.Unmarshal([]byte(b), cfg)
	if err != nil {
		t.Fatalf("unexpected unmarshal error: %s", err)
	}
	if len(cfg.Attrs) != 2 {
		t.Fatalf("unexpected attrs length: %d", len(cfg.Attrs))
	}
	if cfg.Attrs["a"] != 1 {
		t.Fatalf("unexpected attr a: %d", cfg.Attrs["a"])
	}
	if cfg.Attrs["b"] != "2" {
		t.Fatalf("unexpected attr b: %s", cfg.Attrs["b"])
	}

}
