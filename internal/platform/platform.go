package platform

import (
	"github.com/aymantoumi/portlens/internal/types"
)

type Scanner interface {
	ScanListening() ([]*types.PortEntry, error)
}
