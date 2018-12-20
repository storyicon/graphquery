package graphquery

import (
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/storyicon/graphquery/kernel"
)

func TestParseFromString(t *testing.T) {
	type args struct {
		document string
		expr     string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
	}{
		{
			name: "test0",
			args: args{
				document: `
                    <html><body>
                        <div class="item">
                            <div class="title">title A</div>
                            <div class="author">author A</div>
                        </div>
                        <div class="item">
                            <div class="title">title B</div>
                            <div class="author">author B</div>
                        </div>
                        <div class="item">
                            <div class="title">title C</div>
                            <div class="author">author C</div>
                        </div>
                    </body></html>
                `,
				expr: strings.Join([]string{
					"item `css(\".item\")` [",
					"    {",
					"        title `css(\".title\")`",
					"        author `css(\".author\")`",
					"    }",
					"]",
				}, "\r\n"),
			},
			wantResponse: `{"data":[{"author":"author A","title":"title A"},{"author":"author B","title":"title B"},{"author":"author C","title":"title C"}],"errors":null}`,
		},
		{
			name: "test1",
			args: args{
				document: `
                    <html><body>
                        <div class="title">Article</div>
                    </body></html>
                `,
				expr: strings.Join([]string{
					"title `css(\".title\")`",
				}, "\r\n"),
			},
			wantResponse: `{"data":"Article","errors":null}`,
		},
		{
			name: "test2",
			args: args{
				document: `
                    <html><body>
                        <div class="title">Article</div>
                        <div class="tags">
                            <div class="tag">tag0</div>
                            <div class="tag">tag1</div>
                            <div class="tag">tag2</div>
                        </div>
                    </body></html>
                `,
				expr: strings.Join([]string{
					"{",
					"    title `css(\".title\")`",
					"    tags `css(\".tag\")` [",
					"        tag `text()`",
					"    ]",
					"}",
				}, "\r\n"),
			},
			wantResponse: `{"data":{"tags":["tag0","tag1","tag2"],"title":"Article"},"errors":null}`,
		},
		{
			name: "test3",
			args: args{
				document: `
                    <html><body>
                        <div class="items">
                            <div class="item">
                                <div class="pos">1.0</div>
                                <div class="pos">1.1</div>
                                <div class="pos">1.2</div>
                            </div>
                            <div class="item">
                                <div class="pos">2.0</div>
                                <div class="pos">2.1</div>
                                <div class="pos">2.2</div>
                            </div>
                            <div class="item">
                                <div class="pos">3.0</div>
                                <div class="pos">3.1</div>
                                <div class="pos">3.2</div>
                            </div>
                        </div>
                    </body></html>
                `,
				expr: strings.Join([]string{
					"items `css(\".item\")` [",
					"    item `css(\".pos\")` [",
					"        pos `text()`",
					"    ]",
					"]",
				}, "\r\n"),
			},
			wantResponse: `{"data":[["1.0","1.1","1.2"],["2.0","2.1","2.2"],["3.0","3.1","3.2"]],"errors":null}`,
		},
		{
			name: "test4",
			args: args{
				document: `
                    <html><body>
                        <div class="items">
                            <div class="item">
                                <div class="pos">1.0</div>
                                <div class="pos">1.1</div>
                                <div class="pos">1.2</div>
                            </div>
                            <div class="item">
                                <div class="pos">2.0</div>
                                <div class="pos">2.1</div>
                                <div class="pos">2.2</div>
                            </div>
                            <div class="item">
                                <div class="pos">3.0</div>
                                <div class="pos">3.1</div>
                                <div class="pos">3.2</div>
                            </div>
                        </div>
                    </body></html>
                `,
				expr: strings.Join([]string{
					"item `css(\".pos\")` [",
					"    pos `text()`",
					"]",
				}, "\r\n"),
			},
			wantResponse: `{"data":["1.0","1.1","1.2","2.0","2.1","2.2","3.0","3.1","3.2"],"errors":null}`,
		},
		{
			name: "test5",
			args: args{
				document: `
                    <html>
                        <body>
                            <a href="01.html">Page 1</a>
                            <a href="02.html">Page 2</a>
                            <a href="03.html">Page 3</a>
                        </body>
                    </html>
                `,
				expr: strings.Join([]string{
					"{",
					"    anchor `css(\"a\")` [",
					"        content `text()`",
					"    ]",
					"}",
				}, "\r\n"),
			},
			wantResponse: `{"data":{"anchor":["Page 1","Page 2","Page 3"]},"errors":null}`,
		},
		{
			name: "test6",
			args: args{
				document: `
                    <html>
                        <body>
                            <a href="01.html">Page 1</a>
                            <a href="02.html">Page 2</a>
                            <a href="03.html">Page 3</a>
                        </body>
                    </html>
                `,
				expr: strings.Join([]string{
					"{",
					"    anchor `css(\"a\")` [",
					"        {",
					"            title `text()`",
					"        }",
					"    ]",
					"}",
				}, "\r\n"),
			},
			wantResponse: `{"data":{"anchor":[{"title":"Page 1"},{"title":"Page 2"},{"title":"Page 3"}]},"errors":null}`,
		},
		{
			name: "test7",
			args: args{
				document: `
                    <html>
                        <body>
                            <a href="01.html">Page 1</a>
                            <a href="02.html">Page 2</a>
                            <a href="03.html">Page 3</a>
                        </body>
                    </html>
                `,
				expr: strings.Join([]string{
					"anchor `css(\"a\")` [{title `text()`;  url `attr(\"href\")`;}]",
				}, "\r\n"),
			},
			wantResponse: `{"data":[{"title":"Page 1","url":"01.html"},{"title":"Page 2","url":"02.html"},{"title":"Page 3","url":"03.html"}],"errors":null}`,
		},
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for _, tt := range tests {
		if gotResponse := ParseFromString(tt.args.document, tt.args.expr); !reflect.DeepEqual(gotResponse.String(), tt.wantResponse) {
			t.Errorf("%q. ParseFromString() = %v, want %v", tt.name, gotResponse.String(), tt.wantResponse)
		}
	}
}

func TestCompile(t *testing.T) {
	type args struct {
		expr []byte
	}
	tests := []struct {
		name    string
		args    args
		want    kernel.Graph
		wantErr string
	}{
		{
			args: args{
				expr: []byte("[{ title `text();trim()` url  `attr(\"href\")` }]"),
			},
			wantErr: "ReadNode: Unexpected character \"[\", error found in #1 byte of ...|[{ title `t|..., bigger context ...|[{ title `text();trim()` url  `attr(\"href\")` }]|... ",
		},
	}
	for _, tt := range tests {
		_, err := Compile(tt.args.expr)
		if err.Error() != tt.wantErr {
			t.Errorf("%q. Compile() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
	}
}

func TestParseFromBytes(t *testing.T) {
	type args struct {
		document string
		expr     []byte
	}
	tests := []struct {
		name         string
		args         args
		wantResponse string
	}{
		{
			name: "test0",
			args: args{
				document: `<a href="1.html">anchor 1</a> <a href="2.html">anchor 2</a> <a href="3.html">anchor 3</a>`,
				expr:     []byte("a `css(\"a\")` [{ title `text();trim()` url  `attr(\"href\")` }]` }]"),
			},
			wantResponse: `{"data":[{"title":"anchor 1","url":"1.html"},{"title":"anchor 2","url":"2.html"},{"title":"anchor 3","url":"3.html"}],"errors":null}`,
		},
	}
	for _, tt := range tests {
		if gotResponse := ParseFromBytes(tt.args.document, tt.args.expr); !reflect.DeepEqual(gotResponse.JSON(), tt.wantResponse) {
			t.Errorf("%q. ParseFromBytes() = %v, want %v", tt.name, gotResponse.JSON(), tt.wantResponse)
		}
	}
}
