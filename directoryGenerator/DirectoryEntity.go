package directoryGenerator

type Directory struct {
	DirectoryItems []DirectoryItem
	Setting        map[string]string
}

type DirectoryItem struct {
	Title          string
	MarkdownFile   string
	DirectoryItems []DirectoryItem
}
