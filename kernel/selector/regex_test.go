package selector

import (
	"reflect"
	"testing"
)

func TestRegexSelection_Text(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		want     string
		wantErr  bool
	}{
		{
			name:     "test0",
			selector: "Charl.*?chulz",
			want:     "Charles M Schulz",
		},
		{
			name:     "test1",
			selector: "<name>(.*?)</name>",
			want:     "Charles M SchulzPeppermint PattySnoopy",
		},
		{
			name:     "test2",
			selector: "<author id=\"(.*?)\">",
			want:     "CMS",
		},
		{
			name:     "test3",
			selector: "<title lang=\"(.*?)\">(.*?)</title>",
			want:     "en",
		},
	}
	selection, _ := NewRegex(DocumentTest)
	for _, tt := range tests {
		got, err := selection.Find(tt.selector)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. RegexSelection.Text() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got.Text(), tt.want) {
			t.Errorf("%q. RegexSelection.Text() = %v, want %v", tt.name, got.Text(), tt.want)
		}
	}
}
