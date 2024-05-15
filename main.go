package main

import (
	"github.com/RellwNote/RellwNote/directoryGenerator"
	"github.com/RellwNote/RellwNote/log"
)

const filePath = "C:\\Users\\jianing.zhang\\Desktop\\study\\RellwNotes\\mds\\SUMMARY.md"

func main() {
	content := directoryGenerator.GetSummaryFileToByte(filePath)
	directory := directoryGenerator.ParseSummaryByte(content)

	log.Infof("%v", directory)
}
