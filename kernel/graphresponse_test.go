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
	"reflect"
	"testing"
)

func TestGraphResponse_MarshalData(t *testing.T) {
	tests := []struct {
		name     string
		response *GraphResponse
		want     string
		wantErr  bool
	}{
		{
			name: "test0",
			response: &GraphResponse{
				Data: map[string]string{
					"title": "Unknown",
					"link":  "/index",
				},
				Errors: Errors{},
			},
			want: `{"link":"/index","title":"Unknown"}`,
		},
	}
	for _, tt := range tests {
		got, err := tt.response.MarshalData()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GraphResponse.MarshalData() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. GraphResponse.MarshalData() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGraphResponse_JSON(t *testing.T) {
	tests := []struct {
		name       string
		response   *GraphResponse
		wantConseq string
	}{
		{
			name: "test0",
			response: &GraphResponse{
				Data: map[string]string{
					"title": "Unknown",
					"link":  "/index",
				},
				Errors: Errors{},
			},
			wantConseq: `{"data":{"link":"/index","title":"Unknown"},"errors":null}`,
		},
		{
			name: "test1",
			response: &GraphResponse{
				Errors: Errors{
					&Error{
						Err: "Can not parse",
					},
				},
			},
			wantConseq: `{"data":null,"errors":["Can not parse"]}`,
		},
	}
	for _, tt := range tests {
		if gotConseq := tt.response.JSON(); gotConseq != tt.wantConseq {
			t.Errorf("%q. GraphResponse.JSON() = %v, want %v", tt.name, gotConseq, tt.wantConseq)
		}
	}
}

func TestGraphResponse_Decode(t *testing.T) {
	type Anchor struct {
		Title string
		Link  string
	}
	type args struct {
		obj Anchor
	}
	tests := []struct {
		name     string
		response *GraphResponse
		args     args
		wantObj  interface{}
	}{
		{
			name: "test0",
			response: &GraphResponse{
				Data: map[string]string{
					"title": "Unknown",
					"link":  "/index",
				},
				Errors: Errors{},
			},
			args: args{
				obj: Anchor{},
			},
			wantObj: Anchor{
				Title: "Unknown",
				Link:  "/index",
			},
		},
	}
	for _, tt := range tests {
		if err := tt.response.Decode(&tt.args.obj); err == nil && !reflect.DeepEqual(tt.args.obj, tt.wantObj) {
			t.Errorf("%q. GraphResponse.Decode() obj = %v, wantObj %v", tt.name, tt.args, tt.wantObj)
		}
	}
}
