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

package graphquery

import (
	"fmt"

	"storyicon.visualstudio.com/graphquery/compiler"
	"storyicon.visualstudio.com/graphquery/kernel"
)

// Compile compiles an expression in the form of a byte array into a parser
func Compile(expr []byte) (*kernel.Graph, error) {
	return compiler.Compile(expr)
}

// MustCompile is used to compile expressions, will panic when compile errors
func MustCompile(expr []byte) *kernel.Graph {
	parser, err := compiler.Compile(expr)
	if err == nil {
		return parser
	}
	panic(err)
}

// ParseFromString will parse documents and expressions directly,
// and errors will be recorded in the Errors of Response
func ParseFromString(document string, expr string) (response *kernel.GraphResponse) {
	return ParseFromBytes(document, []byte(expr))
}

// ParseFromBytes is used to parse documents in string format and expressions in []byte form
// errors will be recorded in the Errors of Response
func ParseFromBytes(document string, expr []byte) (response *kernel.GraphResponse) {
	parser, err := Compile(expr)
	response = &kernel.GraphResponse{}
	if err != nil {
		response.AddError(fmt.Sprintf("--- Compile Error: %s", err))
		return
	}
	return parser.Parse(document)
}
