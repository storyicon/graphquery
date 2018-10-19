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

import "github.com/storyicon/graphquery/kernel"

// ReadArray attempts to read a array from the byte stream.
func (iter *Iterator) ReadArray() (nodes []*kernel.GraphNode) {
	c := iter.nextToken()
	if c != '[' {
		iter.ReportMismatchChar("ReadStringTuple", '[', c)
	}
	for {
		c = iter.nextToken()
		switch c {
		case ']':
			return
		default:
			iter.unreadByte()
			nodes = append(nodes, iter.ReadNode())
		}
	}
}

// ReadObjectArray attempts to read a object array from the byte stream.
func (iter *Iterator) ReadObjectArray() (nodes []*kernel.GraphNode) {
	c := iter.nextToken()
	if c != '[' {
		iter.ReportMismatchChar("ReadObjectArray", '[', c)
	} else if c = iter.nextToken(); c != '{' {
		iter.ReportMismatchChar("ReadObjectArray", '{', c)
	}
	for {
		c = iter.nextToken()
		switch c {
		case '}':
			head := iter.head
			c = iter.nextToken()
			if c != ']' {
				iter.head = head
				iter.ReportMismatchChar("ReadObjectArray", ']', c)
			}
			return
		default:
			iter.unreadByte()
			nodes = append(nodes, iter.ReadNode())
		}
	}
}
