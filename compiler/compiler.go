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

package compiler

import (
	"fmt"

	"github.com/storyicon/graphquery/kernel"
)

// Iterator is a io.Reader like object, with specific read functions.
type Iterator struct {
	bytes      []byte
	head, tail int
	Error      error
}

// ParseBytes creates an Iterator instance from byte array
func ParseBytes(bytes []byte) *Iterator {
	return &Iterator{
		bytes: bytes,
		head:  0,
		tail:  len(bytes),
	}
}

//nextToken used to return the next non-whitespace token
func (iter *Iterator) nextToken() byte {
	for i := iter.head; i < iter.tail; i++ {
		// a variation of skip whitespaces
		c := iter.bytes[i]
		switch c {
		case ' ', '\n', '\t', '\r':
			continue
		}
		iter.head = i + 1
		return c
	}
	return 0
}

func (iter *Iterator) unreadByte() {
	if iter.Error != nil {
		return
	}
	iter.head--
	return
}

// ReportError record a error in iterator instance with current position.
func (iter *Iterator) ReportError(operation string, msg string) {
	if iter.Error != nil {
		return
	}
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	peekEnd := iter.head + 10
	if peekEnd > iter.tail {
		peekEnd = iter.tail
	}
	parsing := string(iter.bytes[peekStart:peekEnd])
	contextStart := iter.head - 50
	if contextStart < 0 {
		contextStart = 0
	}
	contextEnd := iter.head + 50
	if contextEnd > iter.tail {
		contextEnd = iter.tail
	}
	context := string(iter.bytes[contextStart:contextEnd])
	iter.Error = fmt.Errorf("%s: %s, error found in #%v byte of ...|%s|..., bigger context ...|%s|... ",
		operation, msg, iter.head-peekStart, parsing, context)
	panic(iter.Error)
}

// ReportMismatchChar report an error that does not match the expected character.
func (iter *Iterator) ReportMismatchChar(operation string, expect byte, appear byte) {
	iter.ReportError(operation,
		fmt.Sprintf(`Expect "%s" character, but %s appears.`,
			string(expect), string(appear),
		),
	)
}

// ReportUnExpectedChar reports an error of "unexpected character".
func (iter *Iterator) ReportUnExpectedChar(operation string, appear byte) {
	iter.ReportError(operation,
		fmt.Sprintf(`Unexpected character "%s"`,
			string(appear),
		),
	)
}

// WhatIsNext gets ValueType of relatively next element
func (iter *Iterator) WhatIsNext() SignalType {
	valueType := signalValue[iter.nextToken()]
	iter.unreadByte()
	return valueType
}

// WhatIsNextByte gets byte of relatively next element
func (iter *Iterator) WhatIsNextByte() byte {
	c := iter.nextToken()
	iter.unreadByte()
	return c
}

// Read attempts to read the parser from the byte stream
func (iter *Iterator) Read() (parser *kernel.Graph) {
	parser = &kernel.Graph{
		Root: &kernel.GraphNode{
			Name: kernel.TypeRootNode,
		},
	}

	defer func() { recover() }()

	valueType := iter.WhatIsNext()
	switch valueType {
	case SignalObject:
		parser.GraphType = kernel.TypeObjectGraph
		parser.Nodes = iter.ReadObject()
	default:
		node := iter.ReadNode()
		parser.GraphType = kernel.TypeAtomGraph
		parser.Nodes = append(parser.Nodes, node)
	}
	return
}

// Compile is used to compile expressions
func Compile(expr []byte) (*kernel.Graph, error) {
	iterator := ParseBytes(expr)
	parser := iterator.Read()
	return parser, iterator.Error
}
