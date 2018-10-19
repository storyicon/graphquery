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

	"github.com/tidwall/gjson"
)

// JSONSelection is an element set maintained by the JSON parser.
type JSONSelection struct {
	// Nodes stores the current element collection.
	Nodes *gjson.Result
}

// NewJSON is used to initialize a JSON Selection from the string
// It's a constructor function.
func NewJSON(document string) (*JSONSelection, error) {
	conseq := gjson.Parse(document)
	return &JSONSelection{
		Nodes: &conseq,
	}, nil
}

// Find is used to find the set of elements
// described by the selector in the current collection of elements
// it returns the current element set when the selector is empty.
// It's a standard method of the selection implementation
func (selection *JSONSelection) Find(selector string) (Selection, error) {
	nodes := selection.Nodes
	if nodes != nil && selector != "" {
		conseq := nodes.Get(selector)
		nodes = &conseq
	}
	return &JSONSelection{
		Nodes: nodes,
	}, nil
}

// Type method is used to convert the current Selection to other types.
// It's a standard method of the selection implementation
func (selection *JSONSelection) Type(typename string) (Selection, error) {
	if typename == TypeJSON {
		return selection, nil
	}
	return NewSelection(typename, selection.String())
}

// Eq is used to return the index element in the current element collection
// the index starts at 0
// It's a standard method of the selection implementation
func (selection *JSONSelection) Eq(index int) (Selection, error) {
	if index < 0 {
		return nil, errors.New("method Eq received less than 0 parameters")
	}
	nodes := selection.Nodes
	if nodes.IsArray() {
		conseq := nodes.Array()
		if y := len(conseq); y > 0 && index < y {
			return &JSONSelection{
				Nodes: &conseq[index],
			}, nil
		}
	}
	return nil, nil
}

// Each is used to traverse the current elements
// It's a standard method of the selection implementation
func (selection *JSONSelection) Each(iterator func(int, Selection) bool) error {
	i := -1
	if selection.Nodes != nil {
		selection.Nodes.ForEach(func(key, value gjson.Result) bool {
			i++
			return iterator(i, &JSONSelection{
				Nodes: &value,
			})
		})
	} else {
		iterator(0, selection)
	}
	return nil
}

// Attr is used to obtain values of specified attributes of an element
// JSONSelection does not currently support the Attr method.
// TODO: Implement the Attr method of JSONSelection
// example: {"element": "<a href=\"01.html\">Index</a>"}
// json("element").attr("href") => 01.html
// It's a standard method of the selection implementation
func (selection *JSONSelection) Attr(attr string) (string, error) {
	return "", errors.New("the JSON selector does not currently support the Attr method")
}

// String method is used to return the string of all elements in the current element collection
// the html/xml tag in this text will not be deleted
// It's a standard method of the selection implementation
func (selection *JSONSelection) String() string {
	if selection.Nodes == nil {
		return ""
	}
	return selection.Nodes.String()
}

// Text method of the JSON selector now only outputs the contents of the element collection as it is.
// TODO: Better way to achieve
// example: {"element": "<a href=\"01.html\">Index</a>"}
// json("element").text() => Index
// It's a standard method of the selection implementation
func (selection *JSONSelection) Text() string {
	if element, err := NewCSS(selection.String()); err == nil {
		return element.Text()
	}
	return selection.String()
}
