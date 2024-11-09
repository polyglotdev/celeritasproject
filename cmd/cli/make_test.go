package main

import (
	"testing"

	celeritas "github.com/polyglotdev/celeritasproject"
	"github.com/stretchr/testify/assert"
)

// TestNewGenerator tests the generator factory
func TestNewGenerator(t *testing.T) {
	cc := &celeritas.Celeritas{
		RootPath: "/test",
	}
	cc.DB.DataType = "postgres"

	tests := []struct {
		name     string
		genType  string
		wantType Generator
		wantErr  bool
	}{
		{
			name:     "migration generator",
			genType:  "migration",
			wantType: &MigrationGenerator{},
			wantErr:  false,
		},
		{
			name:     "unknown generator",
			genType:  "unknown",
			wantType: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(innerT *testing.T) {
			got, err := newGenerator(tt.genType, cc)

			if tt.wantErr {
				assert.Error(innerT, err)
				assert.Nil(innerT, got)
			} else {
				assert.NoError(innerT, err)
				assert.IsType(innerT, tt.wantType, got)
			}
		})
	}
}
