function loaded() {
    class IndexMarkdownConvert extends RellwNoteMarkdownConvert{
        PostProcess_Anchor(a) {
            a.href = LinkHrefFilter(a.getAttribute('href'))
            const href = a.getAttribute("href")
            if(href.startsWith("?md="))
                a.setAttribute("href","content.html" + href)
        }
    }

    const contentTag = document.querySelector("article.markdown")
    LoadMarkdownSource("/INDEX.md").then(markdown => {
        const convert = new IndexMarkdownConvert()
        const target = document.querySelector("article.markdown")
        convert.Convert(markdown, target)
    }).catch(err => {
        console.log(err)
    }).finally(() => {
    })
}