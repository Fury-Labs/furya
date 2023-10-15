package types

import (
	"github.com/fury-labs/furya/osmomath"
)

// ConsensusMinFee is a governance set parameter from prop 354 (https://www.mintscan.io/furya/proposals/354)
// Its intended to be .0025 uosmo / gas
var ConsensusMinFee osmomath.Dec = osmomath.NewDecWithPrec(25, 4)
