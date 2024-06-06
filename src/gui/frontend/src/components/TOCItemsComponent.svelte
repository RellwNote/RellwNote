<script lang="ts">
    import type {TOCItem} from "./TOCItem";

    export let item: TOCItem;
    console.log(item)

</script>

<div class="item">
    <div class="itemTitle">
        {item.Title}
    </div>
    <div class="itemFile">
        {item.MarkdownFile}
    </div>
    <div class="operator">
        <span>↑</span>
        <span>↓</span>
    </div>
</div>
<div class="itemChildList">
    {#if item.TOCItems.length !== 0}
        {#each item.TOCItems as itemInner}
            <svelte:self item = {itemInner}></svelte:self>
        {/each}
    {/if}
</div>

<style>
    .item {
        height: 2rem;
        display: flex;
        justify-content: start;
        border: solid 1px #369fb6;
        align-items: center;
    }
    .itemChildList {
        margin-left: 1.5rem;
    }
    .itemTitle {
        margin-right: 1.5rem;
        margin-left: .5rem;
    }
    .operator {
        justify-self: end;
    }

    .operator span {
        width: 2rem;
        height: 2rem;
    }
</style>

