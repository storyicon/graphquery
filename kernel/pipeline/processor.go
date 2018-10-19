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
	"fmt"

	"github.com/storyicon/graphquery/kernel/selector"
)

// Processor is the function call unit of pipeline.
type Processor struct {
	// Func is the function ontology.
	Func Callee
	// ArgsCount is the number of function parameters.
	ArgsCount int
}

// Callee defines the function body of Processor.
type Callee func(selector.Selection, []string) (selector.Selection, error)

const (
	// ErrUndefinedMethod means method undefined
	ErrUndefinedMethod = "undefined method: %s"
	// ErrWrongArgNumber means wrong args number
	ErrWrongArgNumber = "method %s expects %d parameters, but %d received"
	// ErrAlreadyExists means processor already exists
	ErrAlreadyExists = "processor regist failed: %s already exists"
)

// _Registry stores all registered Processor.
var _Registry = map[string]*Processor{}

func getProcessor(name string) *Processor {
	if proc, exists := _Registry[name]; exists {
		return proc
	}
	return nil
}

// RegistProcessor is used to register a Processor with the registry.
func RegistProcessor(name string, callee Callee, argsCount int) error {
	if getProcessor(name) != nil {
		return fmt.Errorf(ErrAlreadyExists, name)
	}
	_Registry[name] = &Processor{
		Func:      callee,
		ArgsCount: argsCount,
	}
	return nil
}

// InvokeProcessor is used to invoke a Processor.
func InvokeProcessor(node selector.Selection, name string, args []string) (selector.Selection, error) {
	proc := getProcessor(name)
	if proc == nil {
		return nil, fmt.Errorf(ErrUndefinedMethod, name)
	}
	if len(args) != proc.ArgsCount {
		return nil, fmt.Errorf(ErrWrongArgNumber,
			name, proc.ArgsCount, len(args),
		)
	}
	return proc.Func(node, args)
}
