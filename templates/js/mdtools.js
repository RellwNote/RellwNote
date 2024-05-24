/**
 * @typedef {Object} HttpError
 * @property {number} state
 * @property {string} message
 */

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