package tmuxcli

import "testing"

func TestParseSpawnResult(t *testing.T) {
	tests := []struct {
		name    string
		out     string
		want    SpawnResult
		wantErr bool
	}{
		{
			name: "created",
			out:  "$1\x1fmain\x1f@2\x1f3\x1f%4\n",
			want: SpawnResult{SessionId: "$1", SessionName: "main", WindowId: "@2", WindowIndex: 3, PaneId: "%4"},
		},
		{name: "nothing created", out: "", want: SpawnResult{}},
		{name: "bad index", out: "$1\x1fmain\x1f@2\x1fx\x1f%4\n", wantErr: true},
		{name: "wrong field count", out: "main:0\n", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSpawnResult(tt.out)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseSpawnResult err = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("parseSpawnResult = %+v, want %+v", got, tt.want)
			}
		})
	}
}
