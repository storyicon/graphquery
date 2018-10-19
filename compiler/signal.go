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

const (
	// SignalInvalid identifies unrecognized signals
	SignalInvalid = iota
	// SignalName identifies the beginning of the function name and key name
	SignalName
	// SignalPipeline identifies the beginning of the pipeline
	SignalPipeline
	// SignalObject identifies the beginning of the object
	SignalObject
	// SignalArray identifies the beginning of the array
	SignalArray
)

// SignalType is the signal value type
type SignalType int

var signalValue []SignalType

func init() {
	signalValue = make([]SignalType, 256)
	for i := '0'; i < '9'; i++ {
		signalValue[i] = SignalName
	}
	for i := 'a'; i < 'z'; i++ {
		signalValue[i] = SignalName
	}
	for i := 'A'; i < 'Z'; i++ {
		signalValue[i] = SignalName
	}
	signalValue['_'] = SignalName
	signalValue['['] = SignalArray
	signalValue['{'] = SignalObject
	signalValue['`'] = SignalPipeline
}
