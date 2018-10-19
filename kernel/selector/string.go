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

import "errors"

// StringSelection is an element set maintained by the string parser.
type StringSelection struct {
	// Nodes stores the current element collection.
	Nodes []string
}

// NewString is used to initialize a String Selection from the string
// It's a constructor function.
func NewString(document string) (Selection, error) {
	return &StringSelection{
		Nodes: []string{
			document,
		},
	}, nil
}

// Find method of StringSelection returns itself.
func (selection *StringSelection) Find(selector string) (Selection, error) {
	return selection, nil
}

// Type method is used to convert the current Selection to other types.
// It's a standard method of the selection implementation
func (selection *StringSelection) Type(typename string) (Selection, error) {
	if typename == TypeSTRING {
		return selection, nil
	}
	return NewSelection(typename, selection.String())
}

// Eq is used to return the index element in the current element collection
// the index starts at 0
// It's a standard method of the selection implementation
func (selection *StringSelection) Eq(index int) (Selection, error) {
	if index < 0 {
		return nil, errors.New("method Eq received less than 0 parameters")
	}
	nodes := selection.Nodes
	if y := len(nodes); y > 0 && index < y {
		return &StringSelection{
			Nodes: []string{
				nodes[index],
			},
		}, nil
	}
	return nil, nil
}

// Each is used to traverse the current elements
// It's a standard method of the selection implementation
func (selection *StringSelection) Each(iterator func(int, Selection) bool) error {
	for i := 0; i < len(selection.Nodes); i++ {
		if !iterator(i, &StringSelection{
			Nodes: []string{
				selection.Nodes[i],
			},
		}) {
			break
		}
	}
	return nil
}

// Attr is used to obtain values of specified attributes of an element
// StringSelection does not currently support the Attr method.
// TODO: Implement the Attr method of StringSelection
// It's a standard method of the selection implementation
func (selection *StringSelection) Attr(attr string) (string, error) {
	return "", nil
}

// String method is used to return the string of all elements in the current element collection
// the html/xml tag in this text will not be deleted
// It's a standard method of the selection implementation
func (selection *StringSelection) String() (document string) {
	for _, node := range selection.Nodes {
		document += node
	}
	return
}

// Text method of the StringSelection now only outputs the contents of the element collection as it is.
// TODO: Implement the Text method of StringSelection
func (selection *StringSelection) Text() string {
	return selection.String()
}
