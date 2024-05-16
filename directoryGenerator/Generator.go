package directoryGenerator

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var mdParser parser.Parser

// init 初始化md parse工具
func init() {
	mdParser = goldmark.DefaultParser()
}

// ParseSummaryByte 传入md数据，转成目录结构
func ParseSummaryByte(content []byte) (directory Directory) {
	reader := text.NewReader(content)
	document := mdParser.Parse(reader)

	var rootDirectoryItem DirectoryItem
	parseList(document, content, &rootDirectoryItem)
	directory.DirectoryItems = rootDirectoryItem.DirectoryItems

	return directory
}

// parseList 传入List节点，生成List目录结构
func parseList(node ast.Node, content []byte, rootDirectoryItem *DirectoryItem) {
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		switch item := c.(type) {
		case *ast.List:
			parseList(item, content, rootDirectoryItem)
		case *ast.ListItem:
			listItemDirectoryItem := &DirectoryItem{
				Title:          getTitle(item, content),
				MarkdownFile:   getLink(item),
				DirectoryItems: make([]DirectoryItem, 0),
			}
			parseList(item, content, listItemDirectoryItem)
			rootDirectoryItem.DirectoryItems = append(rootDirectoryItem.DirectoryItems, *listItemDirectoryItem)
		case *ast.ThematicBreak:
			listItemDirectoryItem := &DirectoryItem{
				Title:          "---",
				DirectoryItems: make([]DirectoryItem, 0),
			}
			rootDirectoryItem.DirectoryItems = append(rootDirectoryItem.DirectoryItems, *listItemDirectoryItem)
		}
	}
}

// getTitle 获取一个节点的标题
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

// getLink 获取一个节点的引用文件
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
