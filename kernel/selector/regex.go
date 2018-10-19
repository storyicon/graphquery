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
	"fmt"
	"regexp"
	"strings"
)

// RegexSelection is an element set maintained by the Regex parser.
type RegexSelection struct {
	// Nodes stores the current element collection.
	Nodes    []string
	selector *RegexSelector
}

// RegexSelector defines a parsed regular selector
type RegexSelector struct {
	// Raw is unprocessed regular selector
	Raw string
	// Selector is a regular expression that golang's regexp package can recognize
	Selector string
	// Regexp is a structured Selector
	Regexp *regexp.Regexp
	// Modifier is a collection of strings of regex modifiers
	Modifier string
	// SubIndex is defined by Modifier. When Modifier contains w, SubIndex will be 1, meaning "without outer"
	// For example, the text is <div><p>hellow</p></div>
	// when the regular expression is /<p>(.*?)</p>/, Modifer is an empty string, does not contain the letter w, then SubIndex is 0, and the result is <p>hellow</p>
	// when the regular expression is /<p>(.*?)</p>/w, Modifer is w, contains the letter w, then SubIndex is 1, and the result is hellow
	SubIndex int
}

// NewRegex is used to initialize a Regex Selection from the string
// It's a constructor function.
func NewRegex(document string) (*RegexSelection, error) {
	return &RegexSelection{
		Nodes: []string{
			document,
		},
	}, nil
}

// HasModifier is used to determine whether the current selector contains the specified modifier
func (selector *RegexSelector) HasModifier(modifier string) bool {
	if index := strings.Index(selector.Modifier, modifier); index == -1 {
		return false
	}
	return true
}

// NewRegexSelector is used to initialize a regular expression parser
func NewRegexSelector(expr string) *RegexSelector {
	tail := strings.LastIndex(expr, "/")
	if tail <= 0 {
		return nil
	}

	r := &RegexSelector{
		Raw:      expr,
		Modifier: expr[tail+1:],
		Selector: expr[1:tail],
	}

	reg, err := regexp.Compile(r.Selector)
	if err != nil {
		return nil
	}

	r.Regexp = reg

	//* Modifier w means submatch without outer
	if r.HasModifier("w") {
		r.SubIndex = 1
	}

	return r
}

// Find is used to find the set of elements
// described by the selector in the current collection of elements
// it returns the current element set when the selector is empty.
// It's a standard method of the selection implementation
func (selection *RegexSelection) Find(selector string) (Selection, error) {
	nodes := selection.Nodes
	expr := NewRegexSelector(selector)
	if expr == nil {
		return selection, fmt.Errorf("Unable to resolve regular expression: %s", selector)
	}
	var conseq []string
	subIndex := expr.SubIndex
	regex := expr.Regexp
	for _, node := range nodes {
		matches := regex.FindAllStringSubmatch(node, -1)
		for _, match := range matches {
			if y := len(match); y > 0 && subIndex < y {
				conseq = append(conseq, match[subIndex])
			}
		}
	}
	return &RegexSelection{
		Nodes: conseq,
	}, nil
}

// Type method is used to convert the current Selection to other types.
// It's a standard method of the selection implementation
func (selection *RegexSelection) Type(typename string) (Selection, error) {
	if typename == TypeREGEX {
		return selection, nil
	}
	return NewSelection(typename, selection.String())
}

// Eq is used to return the index element in the current element collection
// the index starts at 0
// It's a standard method of the selection implementation
func (selection *RegexSelection) Eq(index int) (Selection, error) {
	if index < 0 {
		return nil, errors.New("method Eq received less than 0 parameters")
	}
	nodes := selection.Nodes
	if y := len(nodes); y > 0 && index < y {
		return &RegexSelection{
			Nodes: []string{
				nodes[index],
			},
		}, nil
	}
	return nil, nil
}

// Each is used to traverse the current elements
// It's a standard method of the selection implementation
func (selection *RegexSelection) Each(iterator func(int, Selection) bool) error {
	for i := 0; i < len(selection.Nodes); i++ {
		if !iterator(i, &RegexSelection{
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
// RegexSelection does not currently support the Attr method.
// TODO: Implement the Attr method of RegexSelection
// It's a standard method of the selection implementation
func (selection *RegexSelection) Attr(attr string) (string, error) {
	return "", errors.New("the Regex selector does not currently support the Attr method")
}

// String method is used to return the string of all elements in the current element collection
// the html/xml tag in this text will not be deleted
// It's a standard method of the selection implementation
func (selection *RegexSelection) String() (document string) {
	for _, node := range selection.Nodes {
		document += node
	}
	return
}

// Text method of the Regex selector now only outputs the contents of the element collection as it is.
// TODO: Implement the Text method of RegexSelection
// It's a standard method of the selection implementation
func (selection *RegexSelection) Text() string {
	return selection.String()
}
