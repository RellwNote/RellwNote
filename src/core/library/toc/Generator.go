package toc

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var mdParser = goldmark.DefaultParser()

func GetTOCFromFile(filepath string) (TOC Item, err error) {
	content, err := getSummaryFileToByte(filepath)
	content = removeEmptyLinesFromFile(content)
	if err != nil {
		return
	}
	TOC = parseSummaryByte(content)
	return
}

func parseSummaryByte(content []byte) (directory Item) {
	reader := text.NewReader(content)
	document := mdParser.Parse(reader)

	var rootDirectoryItem Item
	parseList(document, content, &rootDirectoryItem)

	return rootDirectoryItem
}

// parseList 传入List节点，生成List目录结构
func parseList(node ast.Node, content []byte, rootDirectoryItem *Item) {
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		switch item := c.(type) {
		case *ast.List:
			parseList(item, content, rootDirectoryItem)
		case *ast.ListItem:
			listItemDirectoryItem := &Item{
				Title:        getTitle(item, content),
				MarkdownFile: getLink(item),
				TOCItems:     make([]Item, 0),
			}
			parseList(item, content, listItemDirectoryItem)
			rootDirectoryItem.TOCItems = append(rootDirectoryItem.TOCItems, *listItemDirectoryItem)
		case *ast.ThematicBreak:
			listItemDirectoryItem := &Item{
				Title:    "---",
				TOCItems: make([]Item, 0),
			}
			rootDirectoryItem.TOCItems = append(rootDirectoryItem.TOCItems, *listItemDirectoryItem)
		}
	}
}

func getTitle(node ast.Node, content []byte) string {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.TextBlock:
			return string(c.Text(content))
		case *ast.Text:
			return string(c.Text(content))
		}
	}
	return ""
}

func getLink(node ast.Node) (link string) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.Link:
			return string(c.Destination)
		case *ast.TextBlock:
			return getLink(c)
		}
	}
	return ""
}
