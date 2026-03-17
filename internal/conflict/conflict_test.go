package conflict

import (
	"testing"

	"github.com/yourusername/portmap/internal/types"
)

func TestSuggest(t *testing.T) {
	pair := &types.ConflictPair{
		Port: 8080,
		A: &types.PortEntry{
			Port:          8080,
			Kind:          types.SourceDocker,
			ContainerName: "web-container",
		},
		B: &types.PortEntry{
			Port:        8080,
			Kind:        types.SourceProcess,
			ProcessName: "nginx",
			PID:         1234,
		},
	}

	s := Suggest(pair)

	if s.Description == "" {
		t.Error("expected non-empty Description")
	}

	if s.Port != 8080 {
		t.Errorf("expected Port 8080, got %d", s.Port)
	}
}
