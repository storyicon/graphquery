package selector

import (
	"reflect"
	"testing"
)

func TestXpathSelection_Find(t *testing.T) {
	type args struct {
		selector string
	}
	tests := []struct {
		name      string
		selection Selection
		args      args
		want      string
		wantErr   bool
	}{
		{
			name: "test0",
			args: args{
				selector: "//book/title",
			},
			want: "Being a Dog Is a Full-Time Job",
		},
		{
			name: "test1",
			args: args{
				selector: "//author/@id",
			},
			want: "CMS",
		},
		{
			name: "test2",
			args: args{
				selector: "$%@#&",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test3",
			args: args{
				selector: "//author/@ids",
			},
			want: "",
		},
	}
	selection, _ := NewXpath(DocumentTest)
	for _, tt := range tests {
		got, err := selection.Find(tt.args.selector)
		if err != nil {
			if !tt.wantErr {
				t.Errorf("%q. XpathSelection.Find() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			continue
		}
		if !reflect.DeepEqual(got.Text(), tt.want) {
			t.Errorf("%q. XpathSelection.Find() = %v, want %v", tt.name, got.Text(), tt.want)
		}
	}
}

func TestXpathSelection_Attr(t *testing.T) {
	type args struct {
		attr     string
		selector string
	}
	tests := []struct {
		name      string
		selection Selection
		args      args
		want      string
		wantErr   bool
	}{
		{
			name: "test0",
			args: args{
				selector: "//book/title",
				attr:     "lang",
			},
			want: "en",
		},
		{
			name: "test1",
			args: args{
				selector: "//author",
				attr:     "!@$@!ids",
			},
			want: "",
		},
	}
	selection, _ := NewXpath(DocumentTest)
	for _, tt := range tests {
		gotSelection, _ := selection.Find(tt.args.selector)
		got, err := gotSelection.Attr(tt.args.attr)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. XpathSelection.Attr() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. XpathSelection.Attr() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestXpathSelection_String(t *testing.T) {
	type args struct {
		selector string
	}
	tests := []struct {
		name         string
		args         args
		wantDocument string
	}{
		{
			name: "test0",
			args: args{
				selector: "//book/title",
			},
			wantDocument: `<title lang="en">Being a Dog Is a Full-Time Job</title>`,
		},
		{
			name: "test1",
			args: args{
				selector: "/html/body/library/comment()",
			},
			wantDocument: "<!-- Great book. -->",
		},
	}
	selection, _ := NewXpath(DocumentTest)
	for _, tt := range tests {
		got, err := selection.Find(tt.args.selector)
		if err != nil {
			continue
		}
		if gotDocument := got.String(); gotDocument != tt.wantDocument {
			t.Errorf("%q. XpathSelection.String() = %v, want %v", tt.name, gotDocument, tt.wantDocument)
		}
	}
}
