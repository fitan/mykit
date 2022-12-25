package main

import "testing"

func TestGen(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				path: "./dev.yaml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := Gen(tt.args.path); (err != nil) != tt.wantErr {
					t.Errorf("Gen() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
