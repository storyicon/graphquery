**We are looking for contributors**! Please check the [ROADMAP](https://github.com/storyicon/graphquery/blob/master/ROADMAP.md) to see how you can help ❤️          

---

❤️ This is a brand new project, let's give it some star and watch to see how it will develop, the open source community is driven by you !

# GraphQuery [![CircleCI](https://circleci.com/gh/storyicon/graphquery/tree/master.svg?style=svg)](https://circleci.com/gh/storyicon/graphquery/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/storyicon/graphquery)](https://goreportcard.com/report/github.com/storyicon/graphquery)  [![Build Status](https://travis-ci.org/storyicon/graphquery.svg?branch=master)](https://travis-ci.org/storyicon/graphquery) [![Coverage Status](https://coveralls.io/repos/github/storyicon/graphquery/badge.svg?branch=master)](https://coveralls.io/github/storyicon/graphquery?branch=master) [![GoDoc](https://godoc.org/github.com/storyicon/graphquery?status.svg)](https://godoc.org/github.com/storyicon/graphquery) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/storyicon/Lobby)
![GraphQuery](https://raw.githubusercontent.com/storyicon/graphquery/master/docs/screenshot/graphquery.png)

 
GraphQuery is a query language and execution engine tied to any backend service. It is `back-end language independent`.    

Related Projects:        
* [GraphQuery-PlayGround](https://github.com/storyicon/graphquery-playground) : Learn and test GraphQuery in an interactive walkthrough       
* [Document](https://github.com/storyicon/graphquery/wiki) : Detailed documentation of GraphQuery
* [GraphQuery-http](https://github.com/storyicon/graphquery-http) : Cross language solution for GraphQuery
## Catalog
- [Overview](#overview)
    - [Language-independent](#language-independent)
    - [Multiple selector syntax support](#multiple-selector-syntax-support)
    - [Complete function](#complete-function)
    - [Clear data structure & Concise grammar](#clear-data-structure--concise-grammar)
- [Getting Started](#getting-started)
    - [1. First example](#1-first-example)
    - [2. Pipeline](#2-pipeline)
- [Install](#install)
    - [1. Golang:](#1-golang)
    - [2. Other language](#2-other-language)

## Overview 
GraphQuery is an easy to use query language, it has built-in `Xpath/CSS/Regex/JSONpath` selectors and enough built-in `text processing functions`.    
The most amazing thing is that you can use the minimalist GraphQuery syntax to get `any data structure you want`.         

### Language-independent
Use GraphQuery to let you unify text parsing logic on any backend language.    
You won't need to find implementations of Xpath/CSS/Regex/JSONpath selectors between different languages ​​and get familiar with their syntax or explore their compatibility.
 
### Multiple selector syntax support
You can use GraphQuery to parse any text and use your skilled selector. GraphQuery currently supports the following selectors:
1. `Jsonpath` for parsing JSON strings
2. `Xpath` and `CSS` for parsing XML/HTML
3. `Regular expressions` for parsing any text.    

You can use these selectors in any combination in GraphQuery. The rich built-in selector provides great flexibility for your parsing.

### Complete function
Graphquery has some built-in text processing functions like `trim`, `template`, `replace`. If you think these functions don't meet your needs, you can register new custom functions in the pipeline.


### Clear data structure & Concise grammar
With GraphQuery, you won't need to look for parsing libraries when parsing text, nor do you need to write complex nesting and traversal. Simple and clear GraphQuery syntax gives you a clear picture of the data structure.      

![compare](https://raw.githubusercontent.com/storyicon/graphquery/master/docs/screenshot/compare.png) 
           
As you can see from the comparison chart above, the syntax of GraphQuery is so simple that even if you contact it for the first time, you can still understand its meaning and get started quickly.            

At the same time, GraphQuery is also very easy to integrate into your backend data system (any backend language), let's continue to look down.                  

## Getting Started 
GraphQuery consists of query language and pipelines. To guide you through each of these components, we've written an example designed to illustrate the various pieces of GraphQuery. This example is not comprehensive, but it is designed to quickly introduce the core concepts of GraphQuery. The premise of the example is that we want to use GraphQuery to query for information about library books.


### 1. First example

```html
<library>
<!-- Great book. -->
<book id="b0836217462" available="true">
    <isbn>0836217462</isbn>
    <title lang="en">Being a Dog Is a Full-Time Job</title>
    <quote>I'd dog paddle the deepest ocean.</quote>
    <author id="CMS">
        <?echo "go rocks"?>
        <name>Charles M Schulz</name>
        <born>1922-11-26</born>
        <dead>2000-02-12</dead>
    </author>
    <character id="PP">
        <name>Peppermint Patty</name>
        <born>1966-08-22</born>
        <qualification>bold, brash and tomboyish</qualification>
    </character>
    <character id="Snoopy">
        <name>Snoopy</name>
        <born>1950-10-04</born>
        <qualification>extroverted beagle</qualification>
    </character>
</book>
</library>
```
Faced with such a text structure, we naturally think of extracting the following data structure from the text :
```
{
    bookID
    title
    isbn
    quote
    language
    author{
        name
        born
        dead
    }
    character [{
        name
        born
        qualification
    }]
}
```
This is perfect, when you know the data structure you want to extract, you have actually succeeded 80%, the above is the data structure we want, we call it DDL (Data Definition Language) for the time being. let's see how GraphQuery does it:
```graphquery
{
    bookID `css("book");attr("id")`
    title `css("title")`
    isbn `xpath("//isbn")`
    quote `css("quote")`
    language `css("title");attr("lang")`
    author `css("author")` {
        name `css("name")`
        born `css("born")`
        dead `css("dead")`
    }
    character `xpath("//character")` [{
        name `css("name")`
        born `css("born")`
        qualification `xpath("qualification")`
    }]
}
```
As you can see, the syntax of GraphQuery adds some strings wrapped in <b>\`</b> to the DDL. These strings wrapped by <b>\`</b> are called `Pipeline`. We will introduce Pipeline later.
Let's first take a look at what data GraphQuery engine returns to us.
```json
{
    "bookID": "b0836217462",
    "title": "Being a Dog Is a Full-Time Job",
    "isbn": "0836217462",
    "quote": "I'd dog paddle the deepest ocean.",
    "language": "en",
    "author": {
        "born": "1922-11-26",
        "dead": "2000-02-12",
        "name": "Charles M Schulz"
    },
    "character": [
        {
            "born": "1966-08-22",
            "name": "Peppermint Patty",
            "qualification": "bold, brash and tomboyish"
        },
        {
            "born": "1950-10-04",
            "name": "Snoopy",
            "qualification": "extroverted beagle"
        }
    ],
}
```
Wow, it's wonderful. Just like what we want.    
We call the above example Example1, now let's have a brief look at what pipeline is.

### 2. Pipeline
A pipeline is a collection of functions that use the parent element text as an entry parameter to execute the functions in the collection in sequence.
For example, the language field in our previous example is defined as follows:
```graphquery
language `css("title");attr("lang")`
```
The `language` is the field name, `css("title"); attr("lang")` is the pipeline. In this pipeline, GraphQuery first uses the CSS selector to find the `title` node from the document, and the title node will be obtained. Pass the obtained node into the attr() function and get its lang attribute. The whole process is as follows:

![language: document->css("title")->attr("lang")->en](https://raw.githubusercontent.com/storyicon/graphquery/master/docs/screenshot/pipeline-getlang.png) 

In Example1, we not only use the css and attr functions, but also xpath(). It is easy to associate, Xpath() is to select elements with the Xpath selector.
The following is a list of the pipeline functions built into the current version of graphquery:    

| pipeline | prototype | example | introduce 
| ------ | ------ | ------ | ----- |
| css | css(CSSSelector) | css("title") | Use CSS selector to select elements | 
| json | json(JSONSelector) | json("title") | Use json path to select elements | 
| xpath | xpath(XpathSelector) |  xpath("//title") |Use Xpath selector to select elements |
| regex | regex(RegexSelector) | regex("<title>(.*?)</title>") | Use Regex selector to select elements |
| trim | trim() | trim() | Clear spaces and line breaks before and after the string|
| template | template(TemplateStr) | template("[{$}]") | Add characters before and after variables|
| attr | attr(AttributeName) | attr("lang") | Extract the property of the current node|
| eq | eq(Index) | eq("0") | Take the nth element in the current node collection|
| string | string() | string() | Extract the current node native string|
| text | text() | text() | Extract the text of the current node|
| link | link(KeyName) | link("title") | Returns the current text of the specified key|
| replace | replace(A, B) | replace("a", "b") | Replace all A in the current node to B|

More detailed introduction to pipeline and function, please go to docs.

## Install 
GraphQuery is currently only native to Golang, but for other languages, it can be invoked as a service.     

### 1. Golang:
```
go get github.com/storyicon/graphquery
```
Create a new go file :
```golang
package main

import (
	"encoding/json"
	"log"

	"github.com/storyicon/graphquery"
)

func main() {
	document := `
        <html>
            <body>
                <a href="01.html">Page 1</a>
                <a href="02.html">Page 2</a>
                <a href="03.html">Page 3</a>
            </body>
        </html>
    `
	expr := "{ anchor `css(\"a\")` [ content `text()` ] }"
	response := graphquery.ParseFromString(document, expr)
	bytes, _ := json.Marshal(response.Data)
	log.Println(string(bytes))
}
```
Run the go file, the output is as follows : 
```
{"anchor":["Page 1","Page 2","Page 3"]}
```

### 2. Other language
We use the HTTP protocol to provide a cross-language solution for developers to query GraphQuery using any back-end language you want to use to access the specified port after starting the service.      

> [GraphQuery-http](https://github.com/storyicon/graphquery-http) : Cross language solution for GraphQuery        

You can also use RPC for communication, but currently you may need to do this yourself, because the RPC project on GraphQuery is still under development.           
At the same time, We welcome the contributors to write native support code for other languages ​​in GraphQuery.      