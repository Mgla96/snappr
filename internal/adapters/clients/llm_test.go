package clients

import "testing"

func TestModelToContextWindow(t *testing.T) {
	type args struct {
		model ModelType
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: string(GPT3_5Turbo0125),
			args: args{model: GPT3_5Turbo0125},
			want: 16385,
		},
		{
			name: string(GPT4_turbo),
			args: args{model: GPT4_turbo},
			want: 128000,
		},
		{
			name: "not found",
			args: args{model: "not-found"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ModelToContextWindow(tt.args.model); got != tt.want {
				t.Errorf("ModelToContextWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}
