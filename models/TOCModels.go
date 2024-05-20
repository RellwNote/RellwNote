package models

type TOC struct {
	TOCItems []TOCItem
	Setting  map[string]string
}

type TOCItem struct {
	Title        string
	MarkdownFile string
	TOCItems     []TOCItem
}
