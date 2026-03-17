package registry

import (
	"testing"
)

func TestLookup(t *testing.T) {
	t.Run("Lookup 5432 returns PostgreSQL", func(t *testing.T) {
		e := Lookup(5432)
		if e == nil {
			t.Fatal("expected non-nil entry")
		}
		if e.Name != "PostgreSQL" {
			t.Errorf("expected Name == 'PostgreSQL', got '%s'", e.Name)
		}
	})

	t.Run("Lookup 99999 returns nil", func(t *testing.T) {
		e := Lookup(99999)
		if e != nil {
			t.Errorf("expected nil, got %v", e)
		}
	})
}

func TestByCategory(t *testing.T) {
	t.Run("ByCategory database returns entries", func(t *testing.T) {
		entries := ByCategory("database")
		if len(entries) == 0 {
			t.Error("expected non-empty slice for database category")
		}
	})
}

func TestIsConventionPort(t *testing.T) {
	t.Run("IsConventionPort 80 returns true", func(t *testing.T) {
		if !IsConventionPort(80) {
			t.Error("expected IsConventionPort(80) to be true")
		}
	})

	t.Run("IsConventionPort 9999 returns false", func(t *testing.T) {
		if IsConventionPort(9999) {
			t.Error("expected IsConventionPort(9999) to be false")
		}
	})
}
