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
	"errors"

	jsoniter "github.com/json-iterator/go"

	"github.com/mitchellh/mapstructure"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// GraphResponse is the analytic result.
type GraphResponse struct {
	Data   GraphRawData `json:"data"`
	Errors Errors       `json:"errors"`
}

func (response *GraphResponse) String() (conseq string) {
	return response.JSON()
}

// JSON is used to convert response to JSON string
func (response *GraphResponse) JSON() (conseq string) {
	if bytes, err := json.Marshal(response); err == nil {
		conseq = string(bytes)
	}
	return
}

// Decode is used to map Response.Data to a given struct
func (response *GraphResponse) Decode(obj interface{}) error {
	if response.Data == nil {
		return errors.New("can not unmarshal nil")
	}
	return mapstructure.Decode(response.Data, obj)
}

// MarshalData is used to convert Response.Data to a JSON string
func (response *GraphResponse) MarshalData() (string, error) {
	if response.Data == nil {
		return "", errors.New("can not marshal nil")
	}
	return json.MarshalToString(response.Data)
}
