package util

import "testing"

func TestValidateAddress(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test URL",
			args: args{str: "https://practicum.yandex.ru/"},
			want: true,
		},
		{
			name: "Test no URL",
			args: args{str: "sdfsdf"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateAddress(tt.args.str); got != tt.want {
				t.Errorf("ValidateAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
