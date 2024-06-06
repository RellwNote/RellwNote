<script lang="ts">
    import {onMount} from 'svelte';
    import {GetTOCFromFile} from "../../wailsjs/go/toc/TOC";
    import type {TOCItem} from "./TOCItem";
    import TOCItemsComponent from "./TOCItemsComponent.svelte";


    let item: TOCItem;
    let isLoading: boolean = true;

    async function GeneratorTOC(){
         // TODO 后续以配置为准
        let itemTemp: TOCItem;
        try {
            itemTemp = await GetTOCFromFile("./test/SUMMARY.md")
        }catch(err){
            console.log(err)
        }
        return itemTemp
    }

    onMount(async () => {
        item = await GeneratorTOC()
        isLoading = false
    });

</script>



<div class="itemList">
    {#if !isLoading}
        {#each item.TOCItems as itemInner}
            <TOCItemsComponent item={itemInner} />
        {/each}
    {:else}
        <p>Loading...</p>
    {/if}
</div>

<style>
    .itemList {
        max-height: 90vh;
        overflow-y: auto;
    }
</style>

