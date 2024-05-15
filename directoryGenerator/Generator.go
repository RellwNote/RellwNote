package directoryGenerator

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var mdParser parser.Parser

// init
// @Description        初始化md parse工具
// @Create             waterIB 2024-05-15 16:24
func init() {
	mdParser = goldmark.DefaultParser()
}

// ParseSummaryByte
// @Description        传入md数据，转成目录结构
// @Create             waterIB 2024-05-15 16:24
// @Param              content []byte md字节数据
// @Return             Directory 目录结构
func ParseSummaryByte(content []byte) (directory Directory) {
	reader := text.NewReader(content)
	document := mdParser.Parse(reader)
	summaryItems := parseDocument(document, content)
	for _, summaryItem := range summaryItems {
		directory.DirectoryItems = append(directory.DirectoryItems, *summaryItem)
	}

	return directory
}

// parseDocument
// @Description        传入目录根节点，生成目录
// @Create             waterIB 2024-05-15 16:24
// @Param              document ast.Node md文档根节点
// @Param              content []byte md字节数据
// @Return             []*DirectoryItem 目录指针数组
func parseDocument(document ast.Node, content []byte) []*DirectoryItem {
	var directoryItems = make([]*DirectoryItem, 0)
	var directoryItem *DirectoryItem

	for child := document.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.Heading:
			isLink, link := getLink(c)
			directoryItem = &DirectoryItem{
				title:          string(c.Text(content)),
				isLink:         isLink,
				markdownFile:   link,
				kind:           c.Kind().String(),
				directoryItems: make([]DirectoryItem, 0),
			}
			directoryItems = append(directoryItems, directoryItem)
		case *ast.List:
			parseList(child, content, directoryItem)
		case *ast.Paragraph:
			isLink, link := getLink(c)
			paragraphItem := DirectoryItem{
				title:          string(c.Text(content)),
				isLink:         isLink,
				markdownFile:   link,
				kind:           c.Kind().String(),
				directoryItems: make([]DirectoryItem, 0),
			}
			directoryItem.directoryItems = append(directoryItem.directoryItems, paragraphItem)
		}
	}

	return directoryItems
}

// parseList
// @Description        传入List节点，生成List目录结构
// @Create             waterIB 2024-05-15 16:24
// @Param              node ast.Node md文档根节点
// @Param              content []byte md字节数据
// @Param              rootDirectoryItem *DirectoryItem 大标题目录节点指针，List结构最终会存入传入的大标题下
func parseList(node ast.Node, content []byte, rootDirectoryItem *DirectoryItem) {
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		switch item := c.(type) {
		case *ast.List:
			parseList(item, content, rootDirectoryItem)
		case *ast.ListItem:
			isLink, link := getLink(item)
			listItemDirectoryItem := &DirectoryItem{
				title:          getTitle(item, content),
				isLink:         isLink,
				markdownFile:   link,
				kind:           item.Kind().String(),
				directoryItems: make([]DirectoryItem, 0),
			}
			parseList(item, content, listItemDirectoryItem)
			rootDirectoryItem.directoryItems = append(rootDirectoryItem.directoryItems, *listItemDirectoryItem)
		}
	}
}

// getTitle
// @Description        获取一个节点的标题
// @Create             waterIB 2024-05-15 16:24
// @Param              node ast.Node md文档根节点
// @Param              content []byte md字节数据
// @Return             string	节点标题
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

// getLink
// @Description        获取一个节点的引用文件
// @Create             waterIB 2024-05-15 16:24
// @Param              node ast.Node md文档根节点
// @Return             bool	是否被引用
// @Return             string	应用文件路径
func getLink(node ast.Node) (isLinking bool, link string) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.Link:
			return true, string(c.Destination)
		}
	}
	return false, string("")
}
