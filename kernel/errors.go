// Copyright 2018 storyicon@foxmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kernel

import (
	"bytes"
	"fmt"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

// Error represents a error's specification.
type Error struct {
	Err string `json:"error"`
}

// Errors is a collection of Error
type Errors []*Error

var _ error = &Error{}

// JSON used to marshal Error to json map
func (msg *Error) JSON() interface{} {
	return msg.Error()
}

// MarshalJSON implements the json.Marshaller interface.
func (msg *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}

// Error implements the error interface.
func (msg Error) Error() string {
	return msg.Err
}

// Errors returns an array will all the error messages.
// Example:
// 		c.Error(errors.New("first"))
// 		c.Error(errors.New("second"))
// 		c.Error(errors.New("third"))
// 		c.Errors.Errors() // == []string{"first", "second", "third"}
func (a Errors) Errors() []string {
	if len(a) == 0 {
		return nil
	}
	errorStrings := make([]string, len(a))
	for i, err := range a {
		errorStrings[i] = err.Error()
	}
	return errorStrings
}

// JSON used to marshal Errors to json string
func (a Errors) JSON() interface{} {
	switch len(a) {
	case 0:
		return nil
	default:
		json := make([]interface{}, len(a))
		for i, err := range a {
			json[i] = err.JSON()
		}
		return json
	}
}

// MarshalJSON implements the json.Marshaller interface.
func (a Errors) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())
}

func (a Errors) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i, msg := range a {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", i+1, msg.Err)
	}
	return buffer.String()
}
