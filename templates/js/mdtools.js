/**
 * 过滤文档中的 Image URL 使其能够直接访问。
 * 因为 MD 是用锚点链接动态加载的，所以 MD 中图片的相对 URL 其实是错误的，
 * 需要调用一下这个方法，以使用当前页面的 MD 地址进行 URL 转换。
 * @param href {string}
 * @return string
 */
function ImageSrcFilter(href) {
    // 站外链接
    if (/^[a-zA-Z]*:\/\//g.test(href))
        return href
    // 站内绝对链接
    if (href.startsWith("/"))
        return href
    // 站内相对链接
    let path = href.replace(/^\.\//, "")
    if (path.startsWith("/") === false)
        path = "/" + path
    return GetCurrentPageMarkdownURLBase() + path
}

/**
 * 过滤文档中的超链接 URL 使其能够直接访问，原因同 ImageSrcFilter 方法。
 * @param href {string}
 * @return string
 */
function LinkHrefFilter(href) {
    // 非 md 文件
    if (href.toLowerCase().endsWith(".md") === false)
        return href
    // 站外链接
    if (/^[a-zA-Z]*:\/\//g.test(href))
        return href
    // 站内绝对链接
    if (href.startsWith("/"))
        return "#" + href
    // 站内相对链接
    let path = href.replace(/^\.\//, "")
    if (path.startsWith("/") === false)
        path = "/" + path
    return "#" + GetCurrentPageMarkdownURLBase() + path
}

/**
 * 获取当前浏览的 MD 的所在目录
 * @returns {string}
 */
function GetCurrentPageMarkdownURLBase(){
    const currentPath = decodeURIComponent(GetCurrentPageMarkdownURL())
    const lastLimiterIndex = currentPath.lastIndexOf("/")
    if (lastLimiterIndex >= 0)
        return currentPath.slice(0, lastLimiterIndex)
    return ""
}

/**
 * 获取当前页面所指向的 Markdown 地址，也就是 URL 中井号后面的地址
 * @returns {string}
 */
function GetCurrentPageMarkdownURL() {
    const lastSharpIndex = document.URL.lastIndexOf("#")
    return lastSharpIndex > 0 ? document.URL.slice(lastSharpIndex + 1) : ""
}

/**
 * 从服务器获取 Markdown 或者纯文本数据
 * @param path {string} 服务器中的Markdown真实路径
 * @returns {Promise<string>}
 * @throws {HttpError}
 */
async function LoadMarkdownSource(path) {
    let f = await fetch(path)
    if (f.ok === false) {
        throw {state: f.status, message: f.statusText}
    }
    let md = await f.text()
    if (md.ok === false) {
        return ""
    }
    return md
}

/**
 * @typedef {Object} HttpError
 * @property {number} state
 * @property {string} message
 */
