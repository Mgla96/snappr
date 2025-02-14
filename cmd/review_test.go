package cmd

import "testing"

func Test_reviewRun(t *testing.T) {
	type args struct {
		commitSHALocal string
		prNumberLocal  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "commit sha is empty",
			args: args{
				commitSHALocal: "",
				prNumberLocal:  1,
			},
			wantErr: true,
		},
		{
			name: "pr number is empty",
			args: args{
				commitSHALocal: "123",
				prNumberLocal:  0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := reviewRun(tt.args.commitSHALocal, tt.args.prNumberLocal); (err != nil) != tt.wantErr {
				t.Errorf("reviewRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
