package TOCGenerator

import (
	"fmt"
	"github.com/RellwNote/RellwNote/log"
	"os"
)

// GetSummaryFileToByte 打开目录文件，获取其字节
func GetSummaryFileToByte(filePath string, fileName string) []byte {
	// 打开目录文件
	fullFilePath := fmt.Sprintf("%s/%s", filePath, fileName)
	_, err := os.Stat(fullFilePath)
	if os.IsNotExist(err) {
		log.Info.Println("目录文件不存在，开始构建目录...")
		// TODO: 构建目录文件
		_ = CreateSummaryFileByFilePath(filePath)
		//WriteContentToFile()
	} else if err != nil {
		log.Error.Println("目录文件打开失败，错误：", err)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Error.Println("读取目录文件失败,", err)
	}
	return content
}

// WriteContentToFile 写入字节数据到覆盖文件中
func WriteContentToFile(filePath string, content []byte) {
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Error.Println("写入文件", filePath, "失败")
	}
}
