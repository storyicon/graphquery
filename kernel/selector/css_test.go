package selector

import (
	"reflect"
	"testing"
)

func TestCSSSelection_Find(t *testing.T) {
	type args struct {
		selector string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test0",
			args: args{
				selector: "isbn",
			},
			want: "0836217462",
		},
		{
			name: "test1",
			args: args{
				selector: "#b0836217462 > character:nth-child(6) > born",
			},
			want: "1950-10-04",
		},
		{
			name: "test2",
			args: args{
				selector: "#CMS > dead",
			},
			want: "2000-02-12",
		},
	}
	selection, _ := NewCSS(DocumentTest)
	for _, tt := range tests {
		got, err := selection.Find(tt.args.selector)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CSSSelection.Find() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got.Text(), tt.want) {
			t.Errorf("%q. CSSSelection.Find() = %v, want %v", tt.name, got.Text(), tt.want)
		}
	}
}

func TestCSSSelection_String(t *testing.T) {
	type args struct {
		selector string
	}
	tests := []struct {
		name         string
		args         args
		selection    Selection
		wantDocument string
	}{
		{
			name: "test0",
			args: args{
				selector: "isbn",
			},
			wantDocument: "<isbn>0836217462</isbn>",
		},
		{
			name: "test1",
			args: args{
				selector: "name",
			},
			wantDocument: "<name>Charles M Schulz</name><name>Peppermint Patty</name><name>Snoopy</name>",
		},
		{
			name: "test2",
			args: args{
				selector: "#b0836217462 > title",
			},
			wantDocument: `<title lang="en">Being a Dog Is a Full-Time Job</title>`,
		},
	}
	selection, _ := NewCSS(DocumentTest)
	for _, tt := range tests {
		got := tt.selection
		if got == nil {
			got, _ = selection.Find(tt.args.selector)
		}
		if gotDocument := got.String(); gotDocument != tt.wantDocument {
			t.Errorf("%q. CSSSelection.String() = %v, want %v", tt.name, gotDocument, tt.wantDocument)
		}
	}
}
