//    Copyright 2018 storyicon@foxmail.com
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
// TODO: URL Absolute & Response Parse
package kernel

import (
	"fmt"
	"strings"

	"storyicon.visualstudio.com/graphquery/kernel/selector"
)

// Graph is a parsed Graph tree
type Graph struct {
	// Root is the virtual root node
	Root *GraphNode
	// Nodes is the user node.
	Nodes []*GraphNode
	// GraphType identifies the output form of Graph.
	GraphType int
	// Data stores analytic data.
	Data   GraphRawData
	Errors []string
}

const (
	// TypeObjectGraph returns raw data
	TypeObjectGraph = iota
	// TypeAtomGraph returns first data
	TypeAtomGraph
)

const (
	// TypeRootNode is the name of the virtual root node
	TypeRootNode = "__ROOT__"
)

// Parse builds the relationship between user nodes and ROOT nodes,
// parse all node data and returns.
func (graph *Graph) Parse(document string) *GraphResponse {
	graph.Root.Selection, _ = selector.NewString(document)
	graph.Root.derivate(graph.Nodes)
	graph.parse()
	graph.bubbleErrors()

	return &GraphResponse{
		Data:   graph.Data,
		Errors: graph.Errors,
	}
}

// parse parse graph all nodes and write to graph.data.
func (graph *Graph) parse() {
	// parse the GraphNode in turn, temporarily stored in the storage
	storage := GraphObject{}
	for _, node := range graph.Nodes {
		storage[node.Name] = node.Parse()
	}

	// judge output type
	switch graph.GraphType {
	case TypeObjectGraph:
		// typeObjectGraph outputs the data as it is
		graph.Data = OutputObject(storage, false)
	case TypeAtomGraph:
		// typeAtomGraph only extracts the node data of the first non-virtual key
		graph.Data = OutputObject(storage, true)
	default:
		graph.addError(fmt.Sprintf(ErrWrongTypeCall,
			"parse", "graph", graph.GraphType,
		))
	}

}

// BubbleErrors collects node errors from the root node of the tree, traversing the entire tree.
func (graph *Graph) bubbleErrors() {
	storage := map[string]int{}
	graph.Root.traverse(func(node *GraphNode) bool {
		if len(node.errors) > 0 {
			for _, err := range node.errors {
				err := fmt.Sprintf("%s: %s", node.Name, err)
				if _, exists := storage[err]; exists {
					storage[err]++
					continue
				}
				storage[err] = 1
			}
		}
		return true
	})
	for err, times := range storage {
		if times > 1 {
			err = fmt.Sprintf("%s (%d times)", err, times)
		}
		graph.addError(err)
	}
}

func (graph *Graph) addError(errs ...interface{}) {
	for _, err := range errs {
		graph.Errors = append(graph.Errors, fmt.Sprint(err))
	}
}

// IsVisualKey determines whether a key is a virtual key based on the key name
func IsVisualKey(name string) bool {
	return strings.HasPrefix(name, "__") && strings.HasSuffix(name, "__")
}
