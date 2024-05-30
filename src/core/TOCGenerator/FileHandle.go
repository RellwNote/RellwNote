package TOCGenerator

import (
	"github.com/RellwNote/RellwNote/log"
	"os"
	"strings"
)

// getSummaryFileToByte 打开目录文件，获取其字节
func getSummaryFileToByte(filePath string) []byte {
	_, err := os.Stat(filePath)
	if err != nil {
		log.Error.Println("目录文件打开失败，错误：", err)
		return nil
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Error.Println("读取目录文件失败,", err)
		return nil
	}
	content = removeEmptyLinesFromFile(content)
	return content
}

// writeContentToFile 写入字节数据到覆盖文件中
func writeContentToFile(filePath string, content []byte) {
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Error.Println("写入文件", filePath, "失败")
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
