package main

import (
	"fmt"
	"github.com/RellwNote/RellwNote/directoryGenerator"
	"github.com/RellwNote/RellwNote/log"
	"github.com/RellwNote/RellwNote/tempServer"
)

func main() {
	filePath := "test/SummaryTest.md"

	content := directoryGenerator.GetSummaryFileToByte(filePath)
	directory := directoryGenerator.ParseSummaryByte(content)
	content = directoryGenerator.ParseDirectoryToByte(directory)
	fmt.Println(string(content))
	log.Infof("%v", directory)
	tempServer.Start()
}
