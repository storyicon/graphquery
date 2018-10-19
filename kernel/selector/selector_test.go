package selector

import (
	"reflect"
	"testing"
)

const DocumentTest = `
    <library>
    <!-- Great book. -->
    <book id="b0836217462" available="true">
        <isbn>0836217462</isbn>
        <title lang="en">Being a Dog Is a Full-Time Job</title>
        <quote>I'd dog paddle the deepest ocean.</quote>
        <author id="CMS">
            <?echo "go rocks"?>
            <name>Charles M Schulz</name>
            <born>1922-11-26</born>
            <dead>2000-02-12</dead>
        </author>
        <character>
            <name>Peppermint Patty</name>
            <born>1966-08-22</born>
            <qualification>bold, brash and tomboyish</qualification>
        </character>
        <character>
            <name>Snoopy</name>
            <born>1950-10-04</born>
            <qualification>extroverted beagle</qualification>
        </character>
    </book>
    </library>
`

const JSONTest = `
    {
        "name": {"first": "Tom", "last": "Anderson"},
        "age":37,
        "children": ["Sara","Alex","Jack"],
        "fav.movie": "Deer Hunter",
        "link": "<a href=\"index.html\">Index</a>",
        "friends": [
            {"first": "Dale", "last": "Murphy", "age": 44},
            {"first": "Roger", "last": "Craig", "age": 68},
            {"first": "Jane", "last": "Murphy", "age": 47}
        ]
    }
`

func TestNewSelection(t *testing.T) {
	type args struct {
		typename string
		document string
	}
	tests := []struct {
		name    string
		args    args
		want    Selection
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := NewSelection(tt.args.typename, tt.args.document)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewSelection() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewSelection() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
