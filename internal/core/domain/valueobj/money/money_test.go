package money_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/money"
)

func TestMoney_New(t *testing.T) {
	tests := []struct {
		name     string
		currency int
		wantErr  bool
	}{
		{"should initialize money vo", 100, false},
		{"should receive an error if get zero value", 0, true},
		{"should receive an error if the max exceeded", 9999999999999999, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := money.New(tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
