/**
 * 过滤文档中的 Image URL 使其能够直接访问。
 * 因为 MD 是用锚点链接动态加载的，所以 MD 中图片的相对 URL 其实是错误的，
 * 需要调用一下这个方法，以使用当前页面的 MD 地址进行 URL 转换。
 * @param href {string}
 * @return string
 */
function ImageSrcFilter(href) {
    switch (CheckLinkType(href)) {
        case "external":
        case "absolute":
            return href
        case "relative":
            let path = href.replace(/^\.\//, "")
            if (path.startsWith("/") === false)
                path = "/" + path
            return GetCurrentPageMarkdownURLBase() + path
    }
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
    switch (CheckLinkType(href)) {
        case "external":
        case "absolute":
            return "?md=" + href
        case "relative":
            let path = href.replace(/^\.\//, "")
            if (path.startsWith("/") === false)
                path = "/" + path
            return "?md=" + GetCurrentPageMarkdownURLBase() + path
    }
}

/**
 * 检测链接的类型
 * @param link {string}
 * @returns {"external"|"absolute"|"relative"}
 */
function CheckLinkType(link) {
    if (/^[a-zA-Z]*:\/\//g.test(link))
        return "external"
    if (link.startsWith("/"))
        return "absolute"
    return "relative"
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
    const params = new URLSearchParams(window.location.search)
    return params.get("md")
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
 * @return {[Map<string,string>,string]}
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

/**
 * 站内的 MD 链接建议始终这个方法做跳转，这样可以不新页面只加载 MD
 * @param tag {HTMLAnchorElement}
 * @param event {Event}
 */
function LinkClick(tag, event) {
    event.preventDefault()
    const href = decodeURIComponent(tag.getAttribute("href"))

    switch (CheckLinkType(href)) {
        case "external":
            window.open(href)
            break
        case "absolute":
        case "relative":
            history.pushState(null, '', href)
            window.dispatchEvent(window.note.onMarkdownLinkChange)
            break
    }
}

/**
 * @type {{onMarkdownLinkChange: Event}}
 */
window.note = {
    onMarkdownLinkChange: new Event("onMarkdownLinkChange"),
}

/**
 * 符合 RellwNote 的 Markdonw 到 HTML 的转换器
 */
class RellwNoteMarkdownConvert {
    /**
     * 创建转换器
     */
    constructor() {
        /**
         * 后处理器列表，Key 是标签的 CSS 选择器，Value 是后处理器
         * @type {Map<string, function(HTMLElement):void>}
         */
        this.postProcessor = new Map()
        this.postProcessor.set("a", this.PostProcess_Anchor)
        this.postProcessor.set("img", this.PostProcess_Img)
        this.postProcessor.set("h1,h2,h3,h4,h5,h6", this.PostProcess_Title)
    }

    /**
     * 转换 Markdown 并将结果应用到 HTML 标签
     * @param source {string}
     * @param targetTag {HTMLElement}
     */
    Convert(source, targetTag) {
        /**
         * Markdown 源代码，解析过程中可能会修改这个代码
         * @type {string}
         */
        this.source = source
        /**
         * 存放 HTML 结果的 HTML 标签
         * @type {HTMLElement}
         */
        this.targetTag = targetTag
        this.ParseFrontMatter()
        this.MarkedParse()
        this.PostProcess()
    }

    /**
     * 尝试解析 MD 头部的配置信息
     */
    ParseFrontMatter() {
        const res = ParseFrontMatterFromSource(this.source)
        /**
         * MD 的头部配置信息
         * @type {Map<string, string>}
         */
        this.frontMatter = res[0]
        this.source = res[1]
    }

    /**
     * MD 转换成 HTML 并应用到 HTML 标签，这里是 marked 库的包装
     */
    MarkedParse() {
        this.targetTag.innerHTML = marked.parse(this.source)
    }

    /**
     * 对某些标签进行后处理
     */
    PostProcess() {
        for (let k of this.postProcessor.keys()) {
            const processor = this.postProcessor.get(k)
            this.targetTag.querySelectorAll(k).forEach(processor)
        }
    }

    PostProcess_Anchor(a) {
        a.href = LinkHrefFilter(a.getAttribute('href'))
        a.addEventListener("click", e => LinkClick(a, e))
    }

    PostProcess_Img(img) {
        img.src = ImageSrcFilter(img.getAttribute('src'))
    }

    PostProcess_Title(title) {
        title.id = "title-" + title.textContent
    }
}