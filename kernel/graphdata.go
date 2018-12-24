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
	"fmt"
	"strconv"
	"strings"
)

// GraphData is the generic memory of Graph.
type GraphData struct {
	Name  string
	Dtype int

	Dstring      string
	Dfloat64     float64
	Dobject      GraphObject
	DobjectArray []GraphObject
	Darray       []*GraphData
}

// GraphRawData used to store JSON structure
type GraphRawData = interface{}

// GraphObject is GraphData organized by map
type GraphObject map[string]*GraphData

const (
	// ErrWrongTypeCall means call method on wrong type
	ErrWrongTypeCall = "call the %s method on %s of wrong type %d"
	// ErrPushOverflow means data push overflow
	ErrPushOverflow = "graph data push overflow"
)

// NewGraphData is used to initialize a GraphData.
func NewGraphData(name string, dataType int) *GraphData {
	return &GraphData{
		Name:  name,
		Dtype: dataType,
	}
}

// Value used to set atomized values for GraphData
func (data *GraphData) Value(val string) (err error) {
	switch data.Dtype {
	case TypeString:
		data.Dstring = val
	case TypeFloat64:
		data.Dfloat64, err = strconv.ParseFloat(val, 64)
	default:
		return fmt.Errorf(ErrWrongTypeCall,
			"value", "graph data", data.Dtype,
		)
	}
	return
}

// Set is used to set new value to Object type GraphData.
func (data *GraphData) Set(key string, val *GraphData) (err error) {
	switch data.Dtype {
	case TypeObject:
		if data.Dobject == nil {
			data.Dobject = GraphObject{}
		}
		data.Dobject[key] = val
		return nil
	default:
		return fmt.Errorf(ErrWrongTypeCall,
			"set", "graph data", data.Dtype,
		)
	}
}

// Push is used to push new value to Array or ObjectArray type GraphData.
func (data *GraphData) Push(index int, key string, val *GraphData) (err error) {
	switch data.Dtype {
	case TypeArray:
		data.Darray = append(data.Darray, val)
	case TypeObjectArray:
		if y := len(data.DobjectArray); index >= y {
			if index > y {
				return fmt.Errorf(ErrPushOverflow)
			}
			data.DobjectArray = append(data.DobjectArray, GraphObject{
				key: val,
			})
		} else {
			data.DobjectArray[index][key] = val
		}
	default:
		return fmt.Errorf(ErrWrongTypeCall,
			"push", "graph data", data.Dtype,
		)
	}
	return nil
}

// Output is used to output the corresponding type of data from GraphData
func (data *GraphData) Output() GraphRawData {
	switch data.Dtype {
	case TypeString:
		return data.Dstring
	case TypeFloat64:
		return data.Dfloat64
	case TypeArray:
		var conseq []GraphRawData
		for _, unit := range data.Darray {
			if strings.HasPrefix(unit.Name, "@") {
				continue
			}
			conseq = append(conseq, unit.Output())
		}
		return conseq
	case TypeObjectArray:
		conseq := []GraphRawData{}
		for _, unit := range data.DobjectArray {
			conseq = append(conseq, OutputObject(unit, false))
		}
		return conseq
	case TypeObject:
		return OutputObject(data.Dobject, false)
	default:
		// impossible, in fact
		return nil
	}
}

// OutputObject is used to traverse the data in output GraphObject
// if atomOnly is true, only one key will be returned
func OutputObject(data GraphObject, atomOnly bool) GraphRawData {
	conseq := map[string]GraphRawData{}
	for _, cell := range data {
		if IsVisualKey(cell.Name) {
			continue
		}
		conseq[cell.Name] = cell.Output()
		if atomOnly {
			return conseq[cell.Name]
		}
	}
	return conseq
}
