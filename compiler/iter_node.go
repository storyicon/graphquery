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
	"storyicon.visualstudio.com/graphquery/kernel"
	"storyicon.visualstudio.com/graphquery/kernel/pipeline"
)

// ReadNode attempts to read a Node from the byte stream.
func (iter *Iterator) ReadNode() (node *kernel.GraphNode) {
	valueType := iter.WhatIsNext()
	if valueType != SignalName {
		iter.ReportUnExpectedChar("ReadNode", iter.nextToken())
		return
	}
	node = &kernel.GraphNode{
		Name:      iter.ReadNodeName(),
		Pipelines: iter.ReadPipelines(),
	}
	node.NodeType, node.Children = iter.ReadChildren()
	return
}

// ReadNodeName attempts to read a continuous node name from the byte stream.
func (iter *Iterator) ReadNodeName() (name string) {
	for i := iter.head; i < iter.tail; i++ {
		c := iter.bytes[i]
		if c == ' ' {
			name = string(iter.bytes[iter.head:i])
			iter.head = i + 1
			return
		} else if signalValue[c] != SignalName {
			iter.ReportUnExpectedChar("ReadNodeName", c)
		}
	}
	return
}

// ReadPipelines attempts to read pipelines from the byte stream
func (iter *Iterator) ReadPipelines() (pipelines pipeline.Pipelines) {
	c := iter.nextToken()
	if c != '`' {
		iter.ReportMismatchChar("ReadPipelines", '`', c)
	}
	for {
		c = iter.nextToken()
		switch c {
		case '`':
			return
		case ';':
			continue
		default:
			iter.unreadByte()
			pipelines = append(pipelines, iter.ReadPipeline())
		}
	}
}

// ReadPipeline attempts to read pipeline from the byte stream
func (iter *Iterator) ReadPipeline() *pipeline.Pipeline {
	return &pipeline.Pipeline{
		Name: iter.ReadFuncName(),
		Args: iter.ReadStringTuple(),
	}
}

// ReadFuncName attempts to read function name from the byte stream
func (iter *Iterator) ReadFuncName() (name string) {
	for i := iter.head; i < iter.tail; i++ {
		c := iter.bytes[i]
		if c == '(' {
			name = string(iter.bytes[iter.head:i])
			iter.head = i
			return
		} else if signalValue[c] != SignalName {
			iter.ReportUnExpectedChar("ReadFuncName", c)
		}
	}
	return
}

// ReadStringTuple attempts to read function arguments from the byte stream
func (iter *Iterator) ReadStringTuple() (args []string) {
	c := iter.nextToken()
	if c != '(' {
		iter.ReportMismatchChar("ReadStringTuple", '(', c)
	}
Loop:
	for {
		c = iter.nextToken()
		switch c {
		case ')':
			return
		case '"':
			for i := iter.head; i < iter.tail; i++ {
				c = iter.bytes[i]
				if c == '"' {
					args = append(args, string(iter.bytes[iter.head:i]))
					iter.head = i + 1
					if c = iter.WhatIsNextByte(); c == ',' {
						iter.head++
					}
					break
				}
			}
		default:
			iter.ReportUnExpectedChar("ReadStringTuple", c)
			break Loop
		}
	}
	return
}

// ReadChildren attempts to read node children from the byte stream
func (iter *Iterator) ReadChildren() (childType int, children []*kernel.GraphNode) {
	c := iter.nextToken()
	switch c {
	case '[':
		head := iter.head - 1
		c2 := iter.nextToken()
		iter.head = head
		if c2 == '{' {
			return kernel.TypeObjectArray, iter.ReadObjectArray()
		}
		return kernel.TypeArray, iter.ReadArray()
	case '{':
		iter.unreadByte()
		return kernel.TypeObject, iter.ReadObject()
	case ';':
		return
	default:
		iter.unreadByte()
		return
	}
}
