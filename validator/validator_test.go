package validator

import (
	"testing"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequestCreateEstate(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		length  int
		wantErr string
	}{
		{"valid", 10, 10, ""},
		{"min bounds", 1, 1, ""},
		{"max bounds", 5000, 5000, ""},
		{"width below min", 0, 10, "width must be between 1 and 5000"},
		{"width above max", 5001, 10, "width must be between 1 and 5000"},
		{"length below min", 10, 0, "length must be between 1 and 5000"},
		{"length above max", 10, 5001, "length must be between 1 and 5000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequestCreateEstate(generated.CreateEstateRequest{Width: tt.width, Length: tt.length})
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequestCreateTree(t *testing.T) {
	validID := uuid.New().String()
	tests := []struct {
		name     string
		estateId string
		x        int
		y        int
		height   int
		wantErr  string
	}{
		{"valid", validID, 10, 10, 15, ""},
		{"min bounds", validID, 1, 1, 1, ""},
		{"max bounds", validID, 5000, 5000, 30, ""},
		{"invalid estate id", "not-a-uuid", 10, 10, 15, "400|must valid estate id"},
		{"empty estate id", "", 10, 10, 15, "400|must valid estate id"},
		{"non-v4 uuid", "00000000-0000-0000-0000-000000000000", 10, 10, 15, "400|must valid estate id"},
		{"x below min", validID, 0, 10, 15, "x must be between 1 and 5000"},
		{"x above max", validID, 5001, 10, 15, "x must be between 1 and 5000"},
		{"y below min", validID, 10, 0, 15, "y must be between 1 and 5000"},
		{"y above max", validID, 10, 5001, 15, "y must be between 1 and 5000"},
		{"height below min", validID, 10, 10, 0, "height must be between 1 and 30"},
		{"height above max", validID, 10, 10, 31, "height must be between 1 and 30"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequestCreateTree(generated.CreateTreeRequest{X: tt.x, Y: tt.y, Height: tt.height}, tt.estateId)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEstateId(t *testing.T) {
	assert.NoError(t, ValidateEstateId(uuid.New().String()))
	assert.EqualError(t, ValidateEstateId("invalid"), "400|must valid estate id")
}
