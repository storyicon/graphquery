package selector

import (
	"testing"
)

func TestJSONSelection_Text(t *testing.T) {
	type args struct {
		selector string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test0",
			args: args{
				selector: "link",
			},
			want: "Index",
		},
		{
			name: "test1",
			args: args{
				selector: "friends.0.first",
			},
			want: "Dale",
		},
	}
	selection, _ := NewJSON(JSONTest)
	for _, tt := range tests {
		gotSelection, err := selection.Find(tt.args.selector)
		if err != nil {
			continue
		}
		if got := gotSelection.Text(); got != tt.want {
			t.Errorf("%q. JSONSelection.Text() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
