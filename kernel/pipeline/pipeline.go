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

package pipeline

import (
	"strings"

	"github.com/storyicon/graphquery/kernel/selector"
)

const (
	// InvokePlaceholder is used for its own replacement.
	// When InvokePlaceholder appears in args, it will be replaced by the string of its own node.
	InvokePlaceholder = string(26) + "{$}" + string(26)
)

// Pipeline structure is the calling unit in pipeline.
type Pipeline struct {
	Name string
	Args []string
}

// Pipelines defines the entire pipeline process of a selection.
type Pipelines []*Pipeline

func invokePlaceholderRender(node selector.Selection, args []string) []string {
	if len(args) == 0 {
		return args
	}
	var cache string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.Index(arg, InvokePlaceholder) != -1 {
			if cache == "" {
				cache = node.String()
			}
			args[i] = strings.Replace(arg, InvokePlaceholder, cache, -1)
		}
	}
	return args
}

// Process performs the entire pipeline process for selection
func Process(selection selector.Selection, pipes Pipelines) (node selector.Selection, err error) {
	node = selection
	for _, pipe := range pipes {
		args := invokePlaceholderRender(node, pipe.Args)
		if node, err = InvokeProcessor(node, pipe.Name, args); err != nil {
			return
		}
	}
	return
}
