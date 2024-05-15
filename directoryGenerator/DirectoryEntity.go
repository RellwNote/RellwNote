package directoryGenerator

type Directory struct {
	DirectoryItems []DirectoryItem
	Setting        map[string]string
}

type DirectoryItem struct {
	title          string
	isLink         bool
	prefix         string
	kind           string
	markdownFile   string
	setting        map[string]string
	directoryItems []DirectoryItem
}
