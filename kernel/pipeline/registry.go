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
	"net/url"
	"strconv"
	"strings"

	"github.com/storyicon/graphquery/kernel/selector"
)

func init() {
	RegistProcessor("css", calleeCSS, 1)
	RegistProcessor("json", calleeJSON, 1)
	RegistProcessor("xpath", calleeXpath, 1)
	RegistProcessor("regex", calleeRegex, 1)
	RegistProcessor("trim", calleeTrim, 0)
	RegistProcessor("template", calleeTemplate, 1)
	RegistProcessor("attr", calleeAttr, 1)
	RegistProcessor("eq", calleeEq, 1)
	RegistProcessor("string", calleeString, 0)
	RegistProcessor("text", calleeText, 0)
	RegistProcessor("link", calleeLink, 1)
	RegistProcessor("replace", calleeReplace, 2)
	RegistProcessor("absolute", calleeAbsolute, 1)
}

func calleeCSS(node selector.Selection, args []string) (selection selector.Selection, err error) {
	expr := args[0]

	if selection, err = node.Type(selector.TypeCSS); err != nil {
		return
	}
	return selection.Find(expr)
}

func calleeJSON(node selector.Selection, args []string) (selection selector.Selection, err error) {
	expr := args[0]

	if selection, err = node.Type(selector.TypeJSON); err != nil {
		return
	}
	return selection.Find(expr)
}

func calleeXpath(node selector.Selection, args []string) (selection selector.Selection, err error) {
	expr := args[0]

	if selection, err = node.Type(selector.TypeXPATH); err != nil {
		return
	}
	return selection.Find(expr)
}

func calleeRegex(node selector.Selection, args []string) (selection selector.Selection, err error) {
	expr := args[0]

	if selection, err = node.Type(selector.TypeREGEX); err != nil {
		return
	}
	return selection.Find(expr)
}

func calleeTrim(node selector.Selection, args []string) (selection selector.Selection, err error) {
	return selector.NewString(strings.TrimSpace(node.String()))
}

func calleeTemplate(node selector.Selection, args []string) (selection selector.Selection, err error) {
	reference := args[0]
	return selector.NewString(reference)
}

func calleeEq(node selector.Selection, args []string) (selection selector.Selection, err error) {
	eq, err := strconv.Atoi(args[0])

	if err != nil {
		return nil, err
	}

	return node.Eq(eq)
}

func calleeString(node selector.Selection, args []string) (selection selector.Selection, err error) {
	return selector.NewString(node.String())
}

func calleeText(node selector.Selection, args []string) (selection selector.Selection, err error) {
	return selector.NewString(node.Text())
}

func calleeAttr(node selector.Selection, args []string) (selection selector.Selection, err error) {
	attr := args[0]

	conseq, err := node.Attr(attr)
	if err != nil {
		return nil, err
	}

	return selector.NewString(conseq)

}

func calleeLink(node selector.Selection, args []string) (selection selector.Selection, err error) {
	reference := args[0]
	return selector.NewString(reference)
}

func calleeReplace(node selector.Selection, args []string) (selection selector.Selection, err error) {
	old, replace := args[0], args[1]
	conseq := strings.Replace(node.String(), old, replace, -1)
	return selector.NewString(conseq)
}

func calleeAbsolute(node selector.Selection, args []string) (selection selector.Selection, err error) {
	raw, parent := node.String(), args[0]
	var parentURL, rawURL *url.URL
	if parentURL, err = url.Parse(parent); err == nil {
		if rawURL, err = url.Parse(raw); err == nil {
			return selector.NewString(parentURL.ResolveReference(rawURL).String())
		}
	}
	return node, err
}
