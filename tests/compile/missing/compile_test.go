package compile_empty

import (
	"testing"

	"github.com/ncruces/go-sqlite3"
)

func TestCompile_empty(t *testing.T) {
	sqlite3.Path = "sqlite3.wasm"
	_, err := sqlite3.Open(":memory:")
	if err == nil {
		t.Error("want error")
	}
}
