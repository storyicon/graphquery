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

	"github.com/PuerkitoBio/goquery"
)

// CSSSelection is an element set maintained by the CSS parser.
type CSSSelection struct {
	// Nodes stores the current element collection.
	Nodes *goquery.Selection
}

// NewCSS is used to initialize a CSS Selection from the string
// It's a constructor function.
func NewCSS(document string) (Selection, error) {
	node, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err == nil {
		return &CSSSelection{
			Nodes: node.Contents(),
		}, nil
	}
	return nil, err
}

// Find is used to find the set of elements
// described by the selector in the current collection of elements
// it returns the current element set when the selector is empty.
// It's a standard method of the selection implementation
func (selection *CSSSelection) Find(selector string) (Selection, error) {
	nodes := selection.Nodes
	if selector != "" {
		nodes = nodes.Find(selector)
	}
	return &CSSSelection{
		Nodes: nodes,
	}, nil
}

// Type method is used to convert the current Selection to other types.
// It's a standard method of the selection implementation
func (selection *CSSSelection) Type(typename string) (Selection, error) {
	if typename == TypeCSS {
		return selection, nil
	}
	return NewSelection(typename, selection.String())
}

// Eq is used to return the index element in the current element collection
// the index starts at 0
// It's a standard method of the selection implementation
func (selection *CSSSelection) Eq(index int) (Selection, error) {
	conseq := selection.Nodes.Eq(index)
	return &CSSSelection{
		Nodes: conseq,
	}, nil
}

// Each is used to traverse the current elements
// It's a standard method of the selection implementation
func (selection *CSSSelection) Each(iterator func(int, Selection) bool) error {
	for i := 0; i < selection.Nodes.Length(); i++ {
		if !iterator(i, &CSSSelection{
			Nodes: selection.Nodes.Eq(i),
		}) {
			break
		}
	}
	return nil
}

// Attr is used to obtain values of specified attributes of an element
// It's a standard method of the selection implementation
func (selection *CSSSelection) Attr(attr string) (string, error) {
	if attr == "" {
		return "", errors.New("attr method requires a string type parameter")
	}
	conseq, _ := selection.Nodes.Attr(attr)
	return conseq, nil
}

// String method is used to return the string of all elements in the current element collection
// the html/xml tag in this text will not be deleted
// It's a standard method of the selection implementation
func (selection *CSSSelection) String() (document string) {
	// if don't traverse here and use goquery. OuterHtml (selection. Nodes) directly,
	// will get only OuterHtml of the first element
	for i := 0; i < selection.Nodes.Length(); i++ {
		block, _ := goquery.OuterHtml(selection.Nodes.Eq(i))
		document += block
	}
	return
}

// Text method is used to return the text of all elements in the current element collection
// the text will not contain html/xml tags and attributes
// It's a standard method of the selection implementation
func (selection *CSSSelection) Text() string {
	return selection.Nodes.Text()
}
