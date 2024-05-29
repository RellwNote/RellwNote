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
        const contentTag = document.querySelector("article#markdown-content .markdown")
        const convert = new RellwNoteMarkdownConvert()
        convert.Convert(source, contentTag)
    }).catch((e) => {
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
window.addEventListener("onMarkdownLinkChange", e => {
    ReloadCurrentMarkdownByURL()
})

ReloadCurrentMarkdownByURL()
