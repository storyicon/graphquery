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

package kernel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/storyicon/graphquery/kernel/pipeline"
	"github.com/storyicon/graphquery/kernel/selector"
)

// GraphNode is the atomic node in Graph.
type GraphNode struct {
	Name       string
	Definition string
	Pipelines  pipeline.Pipelines
	Children   []*GraphNode
	NodeType   int

	Selection selector.Selection
	Parent    *GraphNode
	errors    []string
}

const (
	//TypeString describes the string type.
	TypeString = iota
	//TypeFloat64 describes the float64 type.
	TypeFloat64
	//TypeObject describes the map[string]GraphData type.
	TypeObject
	//TypeObjectArray describes the []map[string]GraphData type.
	TypeObjectArray
	//TypeArray describes the []string, []float64, [][...] type.
	TypeArray
)

const (
	// ErrWrongArgNumber means wrong number of parameters
	ErrWrongArgNumber = "method %s expects %d parameters, but %d received"
	// ErrFatalError means fatal error occurred
	ErrFatalError = "fatal error occurred while %s"
)

// Parse recursively resolves all nodes
func (node *GraphNode) Parse() *GraphData {
	conseq := NewGraphData(node.Name, node.NodeType)
	selection := node.getSelection()

	if selection != nil {
		switch node.NodeType {
		case TypeString, TypeFloat64:
			if err := conseq.Value(node.String()); err != nil {
				node.addError(err)
			}
		case TypeArray, TypeObjectArray:
			selection.Each(func(i int, element selector.Selection) bool {
				// Clearing cache when the node context changes
				defer (func(node *GraphNode) {
					node.each(func(j int, child *GraphNode) bool {
						child.Selection = nil
						return true
					})
				})(node)

				node.Selection = element
				node.each(func(j int, child *GraphNode) bool {
					child.Parent = node
					if err := conseq.Push(i, child.Name, child.Parse()); err != nil {
						node.addError(err)
						return false
					}
					return true
				})
				return true
			})
		case TypeObject:
			node.each(func(j int, child *GraphNode) bool {
				// Clearing cache when the node context changes
				defer (func(node *GraphNode) {
					child.Selection = nil
				})(node)
				child.Parent = node
				if err := conseq.Set(child.Name, child.Parse()); err != nil {
					node.addError(err)
					return false
				}
				return true
			})
		default:
			node.addError(fmt.Sprintf(ErrWrongTypeCall,
				"parse", "graph node", node.NodeType,
			))
		}
	}

	return conseq
}

// Traverse traverses the subtree from the current node
// traversal contains the current node
func (node *GraphNode) traverse(callback func(node *GraphNode) bool) {
	if !callback(node) {
		return
	}
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			child.traverse(callback)
		}
	}
}

//derivate function build relationships for parent and child elements
func (node *GraphNode) derivate(children []*GraphNode) {
	node.Children = children
	for _, child := range children {
		child.Parent = node
	}
}

func (node *GraphNode) addError(err interface{}) {
	node.errors = append(node.errors, fmt.Sprint(err))
}

// getSelection is used to get the selection of the current node
// from cache or to calculate selection by pipeline.
func (node *GraphNode) getSelection() selector.Selection {
	parent := node.Parent
	//* parent == nil occurs only when node is a root node or a fatal error occurs.
	//if it is a root node, because the selection of the root node is constant, it returns directly.
	if parent == nil {
		if node.Name != TypeRootNode {
			node.addError(fmt.Sprintf(ErrFatalError,
				"get selection",
			))
		}
		return node.Selection
	}
	//When the node selection already exists, it returns directly.
	//* Note that this will skip the node's calculation, so be sure to clear the node's selection value when the node's context changes
	if node.Selection != nil {
		return node.Selection
	}
	//calculate the selection of the current node through pipeline
	conseq, err := pipeline.Process(parent.getSelection(), node.getPipelines())
	if err != nil {
		node.addError(err)
	}
	node.Selection = conseq

	return conseq
}

// getPipelines is used to copy the pipelines of node and render it
func (node *GraphNode) getPipelines() pipeline.Pipelines {
	var pipelines pipeline.Pipelines
	for _, pipe := range node.Pipelines {
		var args []string

		switch pipe.Name {
		// link pipeline has some particularities,
		// it can refer to variables directly instead of {$variable} in strings
		case "link":
			if len(pipe.Args) != 1 {
				node.addError(fmt.Sprintf(ErrWrongArgNumber,
					"link", 1, len(pipe.Args),
				))
				continue
			}
			reference := ""
			varname := pipe.Args[0]
			if varname == node.Name {
				continue
			}
			if conseq := node.lookUp(varname); conseq != nil {
				reference = conseq.String()
			}
			args = append(args, reference)

		default:
			// other pipeline
			for _, arg := range pipe.Args {
				args = append(args, node.render(arg))
			}
		}

		// A copy of a variable rendered pipeline.
		pipelines = append(pipelines, &pipeline.Pipeline{
			Name: pipe.Name,
			Args: args,
		})
	}
	return pipelines
}

//render function replaces the magic variable in the passed string with variable value
func (node *GraphNode) render(s string) string {
	//TODO: evaluate without regexp
	expr := regexp.MustCompile(`{\$(.*?)}`)
	matches := expr.FindAllStringSubmatch(s, -1)

	offset := 0

	for _, match := range matches {
		//in fact, it is impossible.
		if len(match) != 2 {
			continue
		}

		template, varname := match[0], match[1]

		var value string
		//reference self
		if varname == "" || varname == node.Name {
			//it will change dynamically, throw it to pipeline.
			value = pipeline.InvokePlaceholder

		} else if conseq := node.lookUp(varname); conseq != nil {
			// referenced to other variables
			value = conseq.String()

		} else {
			// failed to find the variable described
			continue
		}

		s = strings.Replace(s, template, value, 1)
		offset += len(value) - len(template)

	}

	return s
}

// lookUp searches for the node from the previous sibling node and all parent nodes.
func (node *GraphNode) lookUp(name string) (conseq *GraphNode) {
	if conseq = node.prev(name); conseq == nil {
		conseq = node.parents(name)
	}
	return
}

// parents is used to locate nodes from parent,
// including the siblings of parent and ancestor.
func (node *GraphNode) parents(name string) *GraphNode {
	if parent := node.Parent; parent != nil {
		if siblings := parent.siblings(name); siblings != nil {
			return siblings
		}
		return parent.parents(name)
	}
	return nil
}

// siblings is used to find brotherhood nodes.
func (node *GraphNode) siblings(name string) *GraphNode {
	if parent := node.Parent; parent != nil {
		for _, child := range parent.Children {
			if child.Name == name {
				return child
			}
		}
	}
	return nil
}

// prev is used to find the sibling nodes before the current node.
func (node *GraphNode) prev(name string) *GraphNode {
	if parent := node.Parent; parent != nil {
		for _, child := range node.Parent.Children {
			switch child.Name {
			case name:
				return child
			case node.Name: // if traverse the current node, return nil directly.
				return nil
			}
		}
	}
	return nil
}

// getFirstChild is used to return the first child node.
func (node *GraphNode) getFirstChild() *GraphNode {
	if len(node.Children) > 0 {
		return node.Children[0]
	}
	return nil
}

// String method of the node returns the Text of selection.
func (node *GraphNode) String() string {
	if selection := node.getSelection(); selection != nil {
		return selection.Text()
	}
	return ""
}

// Each is used to traverse all the child nodes of the current node.
func (node *GraphNode) each(iterator func(int, *GraphNode) bool) {
	for i := 0; i < len(node.Children); i++ {
		if !iterator(i, node.Children[i]) {
			break
		}
	}
}
