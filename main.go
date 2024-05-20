package main

import (
	"github.com/RellwNote/RellwNote/TOCGenerator"
	"github.com/RellwNote/RellwNote/config"
)

func main() {
	filePath := config.Config.LibraryPath

	//content := directoryGenerator.GetSummaryFileToByte(filePath, config.SummaryFileName)
	//directory := directoryGenerator.ParseSummaryByte(content)
	//content = directoryGenerator.ParseDirectoryToByte(directory)
	//fmt.Println(string(content))
	//log.Infof("%v", directory)
	toc := TOCGenerator.CreateSummaryFileByFilePath(filePath)
	content := TOCGenerator.ParseDirectoryToByte(toc)
	TOCGenerator.WriteContentToFile("./test.md", content)

	//tempServer.Start()
}
