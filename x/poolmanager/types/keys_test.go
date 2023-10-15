package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fury-labs/furya/v20/x/poolmanager/types"
)

func TestFormatDenomTradePairKey(t *testing.T) {
	tests := map[string]struct {
		denom0      string
		denom1      string
		expectedKey string
	}{
		"happy path": {
			denom0:      "ufury",
			denom1:      "uion",
			expectedKey: "\x04|uion|ufury",
		},
		"reversed denoms get reordered": {
			denom0:      "uion",
			denom1:      "ufury",
			expectedKey: "\x04|uion|ufury",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			formatDenomTradePairKey := types.FormatDenomTradePairKey(tc.denom0, tc.denom1)
			stringFormatDenomTradePairKeyString := string(formatDenomTradePairKey)
			require.Equal(t, tc.expectedKey, stringFormatDenomTradePairKeyString)
		})
	}
}
