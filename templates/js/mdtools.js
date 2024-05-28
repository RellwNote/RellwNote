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
function GetCurrentPageMarkdownURLBase() {
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

/**
 * 解析 MD 文档的前缀配置项目
 * @param matter {string}
 * @return {Map<string,any>}
 */
function ParseFrontMatter(matter) {
    let res = new Map()

    matter.split("\n")
        .map(e => e.trim())
        .filter(e => e.indexOf(":") >= 0)
        .map(e => {
            let limiterIndex = e.indexOf(":")
            let prefix = e.slice(0, limiterIndex).trim()
            let value = e.slice(limiterIndex + 1).trim()
            if (value.startsWith('"') && value.endsWith('"') || value.startsWith("'") && value.endsWith("'")) {
                value = value.slice(1, value.length - 1)
            }
            return [prefix, value]
        })
        .forEach(e => res.set(e[0], e[1]))

    return res;
}

/**
 * 解析 MD 文档的前缀配置项目，并返回去除掉配置项目后的 MD 源码
 * @param source {string}
 * @return {[Map<string,any>,string]}
 */
function ParseFrontMatterFromSource(source) {
    let resMap = new Map()
    if (source.trimStart().startsWith("---")) {
        let readyParseFrontMatter = source.trimStart().slice(3)
        let endIndex = readyParseFrontMatter.indexOf("---")
        if (endIndex >= 0) {
            let matter = readyParseFrontMatter.slice(0, endIndex).trim()
            resMap = ParseFrontMatter(matter)

            // 去除前缀配置，防止被显示到文档中
            source = source.trimStart().slice(3)
            source = source.slice(source.indexOf("---") + 3).trimStart()
        }
    }
    return [resMap, source]
}