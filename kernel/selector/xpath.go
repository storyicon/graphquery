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

package selector

import (
	"errors"
	"strings"

	xmlpath "github.com/storyicon/graphquery/kernel/selector/xmlpath.v2"
)

// XpathSelection is an element set maintained by the Xpath parser.
type XpathSelection struct {
	// Nodes stores the current element collection.
	Nodes []*xmlpath.Node
}

// NewXpath is used to initialize a Xpath Selection from the string
// It's a constructor function.
func NewXpath(document string) (*XpathSelection, error) {
	node, err := xmlpath.ParseHTML(strings.NewReader(document))
	if err != nil {
		return nil, err
	}
	return &XpathSelection{
		Nodes: []*xmlpath.Node{
			node,
		},
	}, nil
}

// Find is used to find the set of elements
// described by the selector in the current collection of elements
// it returns the current element set when the selector is empty.
// It's a standard method of the selection implementation
func (selection *XpathSelection) Find(selector string) (Selection, error) {
	nodes := selection.Nodes
	if selector != "" {
		var conseq []*xmlpath.Node
		expr, err := xmlpath.Compile(selector)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			iter := expr.Iter(node)
			for iter.Next() {
				conseq = append(conseq, iter.Node())
			}
		}
		nodes = conseq
	}
	return &XpathSelection{
		Nodes: nodes,
	}, nil
}

// Type method is used to convert the current Selection to other types.
// It's a standard method of the selection implementation
func (selection *XpathSelection) Type(typename string) (Selection, error) {
	if typename == TypeXPATH {
		return selection, nil
	}
	return NewSelection(typename, selection.String())
}

// Eq is used to return the index element in the current element collection
// the index starts at 0
// It's a standard method of the selection implementation
func (selection *XpathSelection) Eq(index int) (Selection, error) {
	if index < 0 {
		return nil, errors.New("method Eq received less than 0 parameters")
	}
	nodes := selection.Nodes
	if y := len(nodes); y > 0 && index < y {
		return &XpathSelection{
			Nodes: []*xmlpath.Node{
				nodes[index],
			},
		}, nil
	}
	return nil, nil
}

// Each is used to traverse the current elements
// It's a standard method of the selection implementation
func (selection *XpathSelection) Each(iterator func(int, Selection) bool) error {
	for i := 0; i < len(selection.Nodes); i++ {
		if !iterator(i, &XpathSelection{
			Nodes: []*xmlpath.Node{
				selection.Nodes[i],
			},
		}) {
			break
		}
	}
	return nil
}

// Attr is used to obtain values of specified attributes of an element
// It's a standard method of the selection implementation
func (selection *XpathSelection) Attr(attr string) (conseq string, err error) {
	if len(selection.Nodes) == 0 {
		return
	}
	if expr, err := xmlpath.Compile("@" + attr); err == nil {
		iter := expr.Iter(selection.Nodes[0])
		if exists := iter.Next(); exists {
			conseq = iter.Node().String()
		}
	}
	return
}

// String method is used to return the string of all elements in the current element collection
// the html/xml tag in this text will not be deleted
// It's a standard method of the selection implementation
func (selection *XpathSelection) String() (document string) {
	for _, node := range selection.Nodes {
		document += node.Html(true)
	}
	return
}

// Text method is used to return the text of all elements in the current element collection
// the text will not contain html/xml tags and attributes
// It's a standard method of the selection implementation
func (selection *XpathSelection) Text() (document string) {
	for _, node := range selection.Nodes {
		document += node.String()
	}
	return
}
