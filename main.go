package main

import (
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/directoryGenerator"
	"github.com/RellwNote/RellwNote/log"
	"github.com/RellwNote/RellwNote/tempServer"
)

func main() {
	filePath := config.GetPublicConfig.Directory.FilePath

	content := directoryGenerator.GetSummaryFileToByte(filePath)
	directory := directoryGenerator.ParseSummaryByte(content)

	log.Infof("%v", directory)
	tempServer.Start()
}
