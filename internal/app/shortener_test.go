package app

import "testing"

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
	}{
		{
			name:    "Test",
			wantLen: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateRandomString(); len(got) != tt.wantLen {
				t.Errorf("GenerateRandomString() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}
