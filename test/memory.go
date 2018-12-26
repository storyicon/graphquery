package main

import (
	"strings"

	"github.com/storyicon/graphquery"

	"gitlab.wallstcn.com/spider/titan-downloader/kernel"
)

func parse() {
	body := kernel.Download("http://www.tbjijian.com/xwlist.do?tag=1").Body
	expr := strings.Replace(`
        a |css("a")| [{
            title |text();trim();template("{$title}");|
            url  |attr("href")|
        }]
	`, "|", "`", -1)
	graphquery.ParseFromString(body, expr)
}

func main() {
	parse()
}
