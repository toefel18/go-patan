package concrete

import (
	"testing"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Error("new store returnes nil")
	}
}