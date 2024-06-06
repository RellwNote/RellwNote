package toc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"io/fs"
	"os"
	fp "path/filepath"
	"strings"
)

type TOC struct {
	ctx context.Context
}

type Item struct {
	Title        string
	MarkdownFile string
	TOCItems     []Item
}

func NewToc() *TOC {
	return &TOC{}
}

func (a *TOC) startup(ctx context.Context) {
	a.ctx = ctx
}

var mdParser = goldmark.DefaultParser()

func (a *TOC) GetTOCFromFile(filepath string) (TOC Item, err error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return Item{}, err
	}
	content = removeEmptyLinesFromFile(content)

	reader := text.NewReader(content)
	document := mdParser.Parse(reader)

	parseList(document, content, &TOC)
	return TOC, nil
}

func (a *TOC) GeneratorTOCByDir(fileDir string) (TOC Item, err error) {
	TOCItem, err := parseFileToTOC(fileDir)
	if err != nil {
		return TOCItem, err
	}
	return TOCItem, err
}

func (a *TOC) SaveTOCToFile(filepath string, TOC Item) error {
	content := parseDirectoryToByte(TOC)
	err := os.WriteFile(filepath, content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func parseDirectoryToByte(TOC Item) []byte {
	content := bytes.NewBuffer([]byte{})
	directoryItems := TOC.TOCItems
	for _, directoryItem := range directoryItems {
		content.Write(parseDirectoryItem(directoryItem, 0))
	}
	return content.Bytes()
}

func parseDirectoryItem(TOCItem Item, layer int) []byte {
	directoryItemByte := bytes.NewBuffer([]byte{})
	for i := 0; i < layer; i++ {
		directoryItemByte.WriteString("\t")
	}
	if len(TOCItem.MarkdownFile) == 0 {
		directoryItemByte.WriteString(fmt.Sprintf("- %s\n", TOCItem.Title))
	} else {
		directoryItemByte.WriteString(fmt.Sprintf("- [%s](%s)\n", TOCItem.Title, TOCItem.MarkdownFile))
	}
	if len(TOCItem.TOCItems) == 0 {
		return directoryItemByte.Bytes()
	}
	for _, v := range TOCItem.TOCItems {
		childBytes := parseDirectoryItem(v, layer+1)
		directoryItemByte.Write(childBytes)
	}
	return directoryItemByte.Bytes()
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

func removeEmptyLinesFromFile(content []byte) []byte {
	lines := strings.Split(string(content), "\n")
	var nonEmptyLines []string

	// 去除空行
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	// 重新组合非空行
	result := strings.Join(nonEmptyLines, "\n")
	return []byte(result)
}

func parseFileToTOC(filepath string) (TOC Item, err error) {
	summaryDir, err := os.Stat(filepath)
	if err != nil {
		return
	}
	if !summaryDir.IsDir() {
		return TOC, errors.New("选中的路径不是目录")
	}

	TOCItem, err := walkDirToCreateTOCItem(filepath)
	return TOCItem, err
}

func walkDirToCreateTOCItem(filepath string) (Item, error) {
	file, err := os.Stat(filepath)
	var TOCItem Item

	if err != nil {
		return TOCItem, err
	}
	TOCItem.Title = file.Name()
	err = fp.Walk(filepath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// TODO 之后需要替换为配置
		if info.Name() == "SUMMARY.md" {
			return nil
		}
		if path == filepath {
			return nil
		}
		if info.IsDir() {
			item, err := walkDirToCreateTOCItem(path)
			if err != nil {
				return err
			}
			TOCItem.TOCItems = append(TOCItem.TOCItems, item)
			return fp.SkipDir
		} else if info.Name() == "index.md" {
			TOCItem.MarkdownFile = convertLink(path)
		} else if strings.ToLower(info.Name()[len(info.Name())-3:len(info.Name())]) == ".md" {
			TOCItem.TOCItems = append(TOCItem.TOCItems, Item{Title: fp.ToSlash(info.Name())[:len(info.Name())-3], MarkdownFile: convertLink(path)})
		}
		return nil
	})
	if err != nil {
		return TOCItem, err
	}
	return TOCItem, err
}

func convertLink(link string) string {
	return strings.ReplaceAll(fp.ToSlash(link), " ", "%20")
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
