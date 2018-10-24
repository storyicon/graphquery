package selector

import (
	"errors"
	"fmt"
	"regexp"
)

// RegexSelection is an element set maintained by the Regex parser.
type RegexSelection struct {
	// Nodes stores the current element collection.
	Nodes []string
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

// Find is used to find the set of elements
// described by the selector in the current collection of elements
// it returns the current element set when the selector is empty.
// It's a standard method of the selection implementation
func (selection *RegexSelection) Find(selector string) (Selection, error) {
	nodes := selection.Nodes
	regex, err := regexp.Compile(selector)
	if err != nil {
		return selection, fmt.Errorf("Unable to resolve regular expression: %s", selector)
	}
	var conseq []string
	for _, node := range nodes {
		matches := regex.FindAllStringSubmatch(node, -1)
		for _, match := range matches {
			subIndex := 0
			if y := len(match); y > 1 {
				subIndex = 1
			}
			conseq = append(conseq, match[subIndex])
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
