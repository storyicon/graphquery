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

const (
	// TypeJSON identifies the currently selected type as JSON
	TypeJSON = "JSON"
	// TypeCSS identifies the currently selected type as CSS
	TypeCSS = "CSS"
	// TypeXPATH identifies the currently selected type as XPATH
	TypeXPATH = "XPATH"
	// TypeREGEX identifies the currently selected type as REGEX
	TypeREGEX = "REGEX"
	// TypeSTRING identifies the currently selected type as STRING
	TypeSTRING = "STRING"
)

// Selection is the selector interface.
type Selection interface {
	// Find is used to find the set of elements
	// described by the selector in the current collection of elements
	// it returns the current element set when the selector is empty.
	Find(string) (Selection, error)

	// Type method is used to convert the current Selection to other types.
	Type(string) (Selection, error)

	// Eq is used to return the index element in the current element collection
	// the index starts at 0
	Eq(int) (Selection, error)

	// Each is used to traverse the current elements
	Each(func(int, Selection) bool) error

	// Attr is used to obtain values of specified attributes of an element
	Attr(string) (string, error)

	// Text method is used to return the text of all elements in the current element collection
	// the text will not contain html/xml tags and attributes
	Text() string

	// String method is used to return the string of all elements in the current element collection
	// the html/xml tag in this text will not be deleted
	String() string
}

// NewSelection is used to initialize the selector of the specified type from the string.
func NewSelection(typename string, document string) (Selection, error) {
	switch typename {
	case TypeCSS:
		return NewCSS(document)
	case TypeJSON:
		return NewJSON(document)
	case TypeREGEX:
		return NewRegex(document)
	case TypeXPATH:
		return NewXpath(document)
	}
	return nil, errors.New("undefined selection type")
}
