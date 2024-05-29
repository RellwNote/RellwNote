function loaded() {
    const contentTag = document.querySelector("article.markdown")
    LoadMarkdownSource("/INDEX.md").then(markdown => {
        ShowMarkdown(markdown, contentTag)
    }).catch(err => {
        console.log(err)
    }).finally(() => {

    })
}