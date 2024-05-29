/**
 * 将 Markdown 显示到页面上
 * @param source {string} Markdown 源码
 */
function ShowMarkdown(source) {
    let parseResult = ParseFrontMatterFromSource(source)
    window.frontMatter = parseResult[0]
    source = parseResult[1]

    const contentTag = document.querySelector("article#markdown-content .markdown")
    contentTag.innerHTML = marked.parse(source)

    for (let a of contentTag.querySelectorAll("a")) {
        a.href = LinkHrefFilter(a.getAttribute('href'))
    }
    for (let image of contentTag.querySelectorAll("img")) {
        image.src = ImageSrcFilter(image.getAttribute('src'))
    }
    for (let h of contentTag.querySelectorAll("h1,h2,h3,h4,h5,h6")) {
        h.id = "title-" + h.textContent
    }

}

/**
 * 根据URL重新加载页面内容
 */
function ReloadCurrentMarkdownByURL() {
    const articleTag = document.querySelector("article#markdown-content")
    const url = GetCurrentPageMarkdownURL()
    articleTag.classList.add('loading')
    articleTag.classList.remove('notfound')

    LoadMarkdownSource(url).then(source => {
        // 防止多次点击目录导致多次请求顺序错乱
        if (url !== GetCurrentPageMarkdownURL())
            return
        ShowMarkdown(source)
    }).catch(() => {
        articleTag.classList.add('notfound')
    }).finally(() => {
        articleTag.classList.remove('loading')
    })
}

/**
 * 滚动页面到锚点位置
 */
function ScrollIntoAnchors() {
    const hash = decodeURIComponent(window.location.hash).slice(1)
    const hashTag = document.getElementById(hash)
    if (hashTag) {
        hashTag.scrollIntoView({
            behavior: 'smooth',
            block: 'start'
        })
    }
}

window.addEventListener("popstate", e => {
    ReloadCurrentMarkdownByURL()
})

ReloadCurrentMarkdownByURL()
