package kernel

import (
	"testing"

	"github.com/storyicon/graphquery/kernel/pipeline"
)

func TestGraph_Parse(t *testing.T) {
	tests := []struct {
		name     string
		document string
		graph    *Graph
		want     string
	}{
		{
			name: "Test0",
			graph: &Graph{
				Root: &GraphNode{
					Name: "__ROOT__",
				},
				GraphType: TypeAtomGraph,
				Nodes: []*GraphNode{
					{
						Name:     "item",
						NodeType: TypeArray,
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "css",
								Args: []string{
									".pos",
								},
							},
						},
						Children: []*GraphNode{
							{
								Name:     "pos",
								NodeType: TypeString,
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "text",
									},
								},
							},
						},
					},
				},
			},
			document: `
                <html><body>
                    <div class="items">
                        <div class="pos">1.0</div>
                        <div class="pos">1.1</div>
                        <div class="pos">1.2</div>
                        <div class="pos">2.0</div>
                        <div class="pos">2.1</div>
                        <div class="pos">2.2</div>
                        <div class="pos">3.0</div>
                        <div class="pos">3.1</div>
                        <div class="pos">3.2</div>
                    </div>
                </body></html>
            `,
			want: `{"data":["1.0","1.1","1.2","2.0","2.1","2.2","3.0","3.1","3.2"],"errors":null}`,
		},
		{
			name: "Test1",
			graph: &Graph{
				Root: &GraphNode{
					Name: "__ROOT__",
				},
				Nodes: []*GraphNode{
					{
						Name:     "title",
						NodeType: TypeString,
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "css",
								Args: []string{
									".title",
								},
							},
						},
					},
					{
						Name: "tags",
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "css",
								Args: []string{
									".tag",
								},
							},
						},
						NodeType: TypeArray,
						Children: []*GraphNode{
							{
								Name: "tag",
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "text",
									},
								},
							},
						},
					},
				},
			},
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
			want: `{"data":{"tags":["tag0","tag1","tag2"],"title":"Article"},"errors":null}`,
		},
		{
			name: "Test2",
			graph: &Graph{
				Root: &GraphNode{
					Name: "__ROOT__",
				},
				GraphType: TypeAtomGraph,
				Nodes: []*GraphNode{
					{
						Name:     "items",
						NodeType: TypeArray,
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "css",
								Args: []string{
									".item",
								},
							},
						},
						Children: []*GraphNode{
							{
								Name:     "item",
								NodeType: TypeArray,
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "css",
										Args: []string{
											".pos",
										},
									},
								},
								Children: []*GraphNode{
									{
										Name:     "pos",
										NodeType: TypeString,
										Pipelines: []*pipeline.Pipeline{
											{
												Name: "text",
											},
										},
									},
								},
							},
						},
					},
				},
			},
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
			want: `{"data":[["1.0","1.1","1.2"],["2.0","2.1","2.2"],["3.0","3.1","3.2"]],"errors":null}`,
		},
		{
			name: "Test3",
			graph: &Graph{
				Root: &GraphNode{
					Name: "__ROOT__",
				},
				Nodes: []*GraphNode{
					{
						Name: "bookname",
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "css",
								Args: []string{
									"book title",
								},
							},
						},
					},
					{
						Name: "exname",
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "regex",
								Args: []string{
									`/<name>(.*?)</name>/w`,
								},
							},
							{
								Name: "eq",
								Args: []string{
									"3",
								},
							},
						},
					},
					{
						Name: "__JSON__",
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "regex",
								Args: []string{
									`/var data = ([\s\S]*?)\s+</script>/w`,
								},
							},
						},
					},
					{
						Name: "__firstname__",
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "link",
								Args: []string{
									"__JSON__",
								},
							},
							{
								Name: "json",
								Args: []string{
									"name.first",
								},
							},
						},
					},
					{
						Name: "author",
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "link",
								Args: []string{
									"__JSON__",
								},
							},
							{
								Name: "json",
								Args: []string{
									"name.last",
								},
							},
							{
								Name: "template",
								Args: []string{
									"{$__firstname__} {$}",
								},
							},
						},
					},
					{
						Name:     "friends",
						NodeType: TypeObjectArray,
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "link",
								Args: []string{
									"__JSON__",
								},
							},
							{
								Name: "json",
								Args: []string{
									"friends",
								},
							},
						},
						Children: []*GraphNode{
							{
								Name: "first",
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "json",
										Args: []string{"first"},
									},
								},
							},
							{
								Name: "last",
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "json",
										Args: []string{"last"},
									},
								},
							},
							{
								Name: "age",
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "json",
										Args: []string{"age"},
									},
								},
							},
						},
					},
					{
						Name:     "character",
						NodeType: TypeObjectArray,
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "xpath",
								Args: []string{"//character"},
							},
						},
						Children: []*GraphNode{
							{
								Name: "name",
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "css",
										Args: []string{"name"},
									},
								},
							},
							{
								Name: "born",
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "xpath",
										Args: []string{"born"},
									},
								},
							},
							{
								Name: "qualification",
								Pipelines: []*pipeline.Pipeline{
									{
										Name: "regex",
										Args: []string{"/qualification>(.*?)</w"},
									},
								},
							},
						},
					},
				},
			},
			document: `
                <!DOCTYPE html>
                <html>
                <head>
                    <meta charset="utf-8" />
                    <meta http-equiv="X-UA-Compatible" content="IE=edge">
                    <title>Test Title</title>
                    <meta name="viewport" content="width=device-width, initial-scale=1">
                    <script>
                        var data = {
                            "name": {
                                "first": "Tom",
                                "last": "Anderson"
                            },
                            "age": 37,
                            "children": ["Sara", "Alex", "Jack"],
                            "fav.movie": "Deer Hunter",
                            "friends": [{
                                    "first": "Dale",
                                    "last": "Murphy",
                                    "age": 44
                                },
                                {
                                    "first": "Roger",
                                    "last": "Craig",
                                    "age": 68
                                },
                                {
                                    "first": "Jane",
                                    "last": "Murphy",
                                    "age": 47
                                }
                            ]
                        }
                    </script>
                </head>
                <body>
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
                            <character id="PP">
                                <name>Peppermint Patty</name>
                                <born>1966-08-22</born>
                                <qualification>bold, brash and tomboyish</qualification>
                            </character>
                            <character id="Snoopy">
                                <name>Snoopy</name>
                                <born>1950-10-04</born>
                                <qualification>extroverted beagle</qualification>
                            </character>
                            <name>Harry Poter</name>
                        </book>
                    </library>
                </body>
                </html>
		    `,
			want: `{"data":{"author":"Tom Anderson","bookname":"Being a Dog Is a Full-Time Job","character":[{"born":"1966-08-22","name":"Peppermint Patty","qualification":"bold, brash and tomboyish"},{"born":"1950-10-04","name":"Snoopy","qualification":"extroverted beagle"}],"exname":"Harry Poter","friends":[{"age":"44","first":"Dale","last":"Murphy"},{"age":"68","first":"Roger","last":"Craig"},{"age":"47","first":"Jane","last":"Murphy"}]},"errors":null}`,
		}, {
			name: "Test4",
			document: `
                <html><body>
                    <div class="title">Article</div>
                </body></html>
            `,
			graph: &Graph{
				Root: &GraphNode{
					Name: "__ROOT__",
				},
				GraphType: TypeAtomGraph,
				Nodes: []*GraphNode{
					{
						Name:     "title",
						NodeType: TypeString,
						Pipelines: []*pipeline.Pipeline{
							{
								Name: "css",
								Args: []string{
									".title",
								},
							},
						},
					},
				},
			},
			want: `{"data":"Article","errors":null}`,
		},
	}

	for _, tt := range tests {
		got := tt.graph.Parse(tt.document).String()
		if tt.want != got {
			t.Errorf("%q. Graph.Parse() = %v, want %v", tt.name, got, tt.want)
		}
	}

}
