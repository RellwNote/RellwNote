note.extensions.push({
    CreateTitle:(type,title) => {
        const p = document.createElement("p")
        p.classList.add("title")

        const iconSpan = document.createElement("span")
        iconSpan.innerText = type.toUpperCase()
        iconSpan.classList.add("icon")
        const titleSpan = document.createElement("span")
        titleSpan.innerText = title
        titleSpan.classList.add("text")
        p.append(iconSpan)
        p.append(titleSpan)
        return p
    },
    MarkdownPostprocessor:function(doc) {
        const blockquotes = doc.querySelectorAll("blockquote")
        for (let block of blockquotes) {
            const title = block.querySelector("p")
            const titleText = title.textContent
            const match = titleText.match(/\[!(.*?)](.*)/)
            if (match === null)
                continue
            const newTitle = this.CreateTitle(match[1],match[2])
            block.classList.add("admonish-obsidian-style")
            block.classList.add("admonition-" + match[1])
            block.removeChild(title)

            const div = document.createElement("div");
            div.classList.add("body")
            while(block.childElementCount > 0) {
                const node = block.childNodes[1]
                block.removeChild(node)
                div.append(node)
            }
            block.append(newTitle)
            block.append(div)
        }
    },
})