package directoryGenerator

import (
	"github.com/RellwNote/RellwNote/log"
	"os"
)

// GetSummaryFileToByte
// @Description        打开目录文件，获取其字节
// @Create             waterIB 2024-05-14 15:56
// @Param              filePath string 文件路径
// @Return             []byte 打开文件的字节
func GetSummaryFileToByte(filePath string) []byte {
	// 打开目录文件
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Info.Println("目录文件不存在，开始构建目录...")
		// TODO: 构建目录文件
	} else if err != nil {
		log.Error.Println("目录文件打开失败，错误：", err)
	}
	content, err := os.ReadFile("C:\\Users\\jianing.zhang\\Desktop\\study\\RellwNotes\\mds\\summary.md")
	if err != nil {
		log.Error.Println("读取目录文件失败,", err)
	}
	return content
}
