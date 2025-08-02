package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentPR(t *testing.T) {
	// Save original state
	originalPRNumber := prNumber
	defer func() {
		prNumber = originalPRNumber
	}()

	tests := []struct {
		name      string
		setPRNum  int
		wantPR    int
		wantErr   bool
		errSubstr string
	}{
		{
			name:     "returns prNumber when set",
			setPRNum: 123,
			wantPR:   123,
			wantErr:  false,
		},
		{
			name:     "returns prNumber when set to positive",
			setPRNum: 456,
			wantPR:   456,
			wantErr:  false,
		},
		// Removed test that calls real gh CLI - covered by integration tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prNumber = tt.setPRNum

			pr, err := getCurrentPR()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errSubstr != "" {
					assert.Contains(t, err.Error(), tt.errSubstr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPR, pr)
			}
		})
	}
}

func TestGetCurrentPREdgeCases(t *testing.T) {
	// Save original state
	originalPRNumber := prNumber
	defer func() {
		prNumber = originalPRNumber
	}()

	// Test edge case where prNumber is negative (should still return it)
	prNumber = -1
	pr, err := getCurrentPR()
	assert.NoError(t, err)
	assert.Equal(t, -1, pr)

	// Test edge case where prNumber is large number
	prNumber = 999999
	pr, err = getCurrentPR()
	assert.NoError(t, err)
	assert.Equal(t, 999999, pr)
}
