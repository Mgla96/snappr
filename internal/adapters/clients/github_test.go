package clients

import "testing"

func TestIsDoNotEditFile(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Contains DO NOT EDIT",
			args: args{
				data: []byte("DO NOT EDIT\nFoobar"),
			},
			want: true,
		},
		{
			name: "Empty file",
			args: args{
				data: []byte(""),
			},
			want: false,
		},
		{
			name: "Partial match",
			args: args{
				data: []byte("DO NO EDIT\nSome other content"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDoNotEditFile(tt.args.data); got != tt.want {
				t.Errorf("IsDoNotEditFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
