package models

type TOCItem struct {
	Title        string
	MarkdownFile string
	TOCItems     []TOCItem
}
