package directoryGenerator

import (
	"bytes"
	"fmt"
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

// ParseDirectoryToByte 将目录结构转化成[]byte存入文件
func ParseDirectoryToByte(directory Directory) []byte {
	content := bytes.NewBuffer([]byte{})
	directoryItems := directory.DirectoryItems
	for _, directoryItem := range directoryItems {
		content.Write(parseDirectoryItem(directoryItem, 0))
	}
	return content.Bytes()
}

// parseDirectoryItem 转化DirectoryItem成[]byte
func parseDirectoryItem(directoryItem DirectoryItem, layer int) []byte {
	directoryItemByte := bytes.NewBuffer([]byte{})
	for i := 0; i < layer; i++ {
		directoryItemByte.WriteString("\t")
	}
	if len(directoryItem.MarkdownFile) == 0 {
		directoryItemByte.WriteString(fmt.Sprintf("- %s\n", directoryItem.Title))
	} else {
		directoryItemByte.WriteString(fmt.Sprintf("- [%s](%s)\n", directoryItem.Title, directoryItem.MarkdownFile))
	}
	if len(directoryItem.DirectoryItems) == 0 {
		return directoryItemByte.Bytes()
	}
	for _, v := range directoryItem.DirectoryItems {
		childBytes := parseDirectoryItem(v, layer+1)
		directoryItemByte.Write(childBytes)
	}
	return directoryItemByte.Bytes()
}
