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
package kernel

import (
	"encoding/json"
	"fmt"
)

// GraphResponse is the analytic result.
type GraphResponse struct {
	Data   GraphRawData `json:"data"`
	Errors []string     `json:"errors"`
}

func (response *GraphResponse) String() (conseq string) {
	if bytes, err := json.Marshal(response); err == nil {
		conseq = string(bytes)
	}
	return
}

func (response *GraphResponse) AddError(err interface{}) {
	response.Errors = append(response.Errors, fmt.Sprint(err))
}
