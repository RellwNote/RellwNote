package toc

import (
	"os"
	"rellwnote/core/log"
	"strings"
)

// getSummaryFileToByte 打开目录文件，获取其字节
func getSummaryFileToByte(filePath string) ([]byte, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// writeContentToFile 写入字节数据到覆盖文件中
func writeContentToFile(filePath string, content []byte) error {
	err := os.WriteFile(filePath, content, os.ModePerm)
	if err != nil {
		log.Error.Println("写入文件", filePath, "失败")
	}
	return err
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
