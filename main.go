package main

import (
	"github.com/RellwNote/RellwNote/TOCGenerator"
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/log"
)

func main() {
	filePath := config.LibraryPath
	summaryName := config.SummaryFileName

	content := TOCGenerator.GetSummaryFileToByte(filePath, summaryName)
	TOCItems := TOCGenerator.ParseSummaryByte(content)
	log.Info.Println(TOCItems)
	//tempServer.Start()
}
