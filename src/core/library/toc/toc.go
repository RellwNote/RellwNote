package toc

type Item struct {
	Title        string
	MarkdownFile string
	TOCItems     []Item
}
