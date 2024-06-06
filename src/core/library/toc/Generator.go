package toc

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"io/fs"
	"os"
	fp "path/filepath"
	"rellwnote/core/config"
	"rellwnote/core/log"
	"rellwnote/core/models"
	"strings"
)

var mdParser = goldmark.DefaultParser()

func GetTOCFromFile(filepath string) (TOC models.TOCItem, err error) {
	content, err := getSummaryFileToByte(filepath)
	TOC = parseSummaryByte(content)
	return
}

func CreateSummaryFileByDirectory(fileDir string, summaryFileName string) error {
	fullFilePath := fp.Join(fileDir, summaryFileName)
	TOCItem, err := parseFileToTOC(fileDir)
	if err != nil {
		return err
	}
	content := parseDirectoryToByte(TOCItem)
	err = writeContentToFile(fullFilePath, content)
	return err
}

// parseSummaryByte 传入md数据，转成目录结构
func parseSummaryByte(content []byte) (directory models.TOCItem) {
	reader := text.NewReader(content)
	document := mdParser.Parse(reader)

	var rootDirectoryItem models.TOCItem
	parseList(document, content, &rootDirectoryItem)

	return rootDirectoryItem
}

// parseList 传入List节点，生成List目录结构
func parseList(node ast.Node, content []byte, rootDirectoryItem *models.TOCItem) {
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		switch item := c.(type) {
		case *ast.List:
			parseList(item, content, rootDirectoryItem)
		case *ast.ListItem:
			listItemDirectoryItem := &models.TOCItem{
				Title:        getTitle(item, content),
				MarkdownFile: getLink(item),
				TOCItems:     make([]models.TOCItem, 0),
			}
			parseList(item, content, listItemDirectoryItem)
			rootDirectoryItem.TOCItems = append(rootDirectoryItem.TOCItems, *listItemDirectoryItem)
		case *ast.ThematicBreak:
			listItemDirectoryItem := &models.TOCItem{
				Title:    "---",
				TOCItems: make([]models.TOCItem, 0),
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

// parseDirectoryToByte 将目录结构转化成[]byte
func parseDirectoryToByte(TOC models.TOCItem) []byte {
	content := bytes.NewBuffer([]byte{})
	directoryItems := TOC.TOCItems
	for _, directoryItem := range directoryItems {
		content.Write(parseDirectoryItem(directoryItem, 0))
	}
	return content.Bytes()
}

func parseDirectoryItem(TOCItem models.TOCItem, layer int) []byte {
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

// parseFileToTOC 通过文件生成目录结构
func parseFileToTOC(filepath string) (TOC models.TOCItem, err error) {
	summaryDir, err := os.Stat(filepath)
	if err != nil {
		log.Error.Println("打开目录文件失败:", err)
	}
	if !summaryDir.IsDir() {
		log.Error.Println("选中的路径：", filepath, "不是目录.", err)
	}

	TOCItem, err := walkDirToCreateTOCItem(filepath)
	return TOCItem, err
}

func walkDirToCreateTOCItem(filepath string) (models.TOCItem, error) {
	file, err := os.Stat(filepath)
	var TOCItem models.TOCItem

	if err != nil {
		log.Error.Println("打开目录文件失败:", err)
		return TOCItem, err
	}
	TOCItem.Title = file.Name()
	err = fp.Walk(filepath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Error.Println("打开文件路径:", path, "失败.", err)
		}
		if info.Name() == config.SummaryFileName {
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
			TOCItem.TOCItems = append(TOCItem.TOCItems, models.TOCItem{Title: fp.ToSlash(info.Name())[:len(info.Name())-3], MarkdownFile: convertLink(path)})
		}
		return nil
	})
	if err != nil {
		log.Error.Println("读取文件夹路径失败: ", filepath, err)
	}
	return TOCItem, err
}

func convertLink(link string) string {
	return strings.ReplaceAll(fp.ToSlash(link), " ", "%20")
}
