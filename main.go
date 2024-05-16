package main

import (
	"github.com/RellwNote/RellwNote/directoryGenerator"
	"github.com/RellwNote/RellwNote/log"
)

const filePath = "./test/SummaryTest.md"

func main() {
	content := directoryGenerator.GetSummaryFileToByte(filePath)
	directory := directoryGenerator.ParseSummaryByte(content)

	log.Infof("%v", directory)
}
