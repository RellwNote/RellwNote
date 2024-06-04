<script lang="ts">
    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';

    interface Item {
    id: number;
    text: string;
}

    let items: Item[] = [
    { id: 1, text: 'Item 1' },
    { id: 2, text: 'Item 2' },
    { id: 3, text: 'Item 3' }
    ];

    function handleDragStart(event: DragEvent) {
    const target = event.target as HTMLElement;
    event.dataTransfer.setData('text/plain', target.id);
}

    function handleDrop(event: DragEvent) {
    const id = event.dataTransfer.getData('text/plain');
    const droppedItem = items.find((item) => item.id == Number(id));
    const newIndex = Number(event.target.dataset.index);

    items = items.filter((item) => item.id != Number(id));
    items = [
    ...items.slice(0, newIndex),
    droppedItem as Item,
    ...items.slice(newIndex)
    ];
}

    onMount(() => {
    const itemsContainer = document.getElementById('items-container');
    if (itemsContainer) {
    itemsContainer.addEventListener('dragover', (event) => {
    event.preventDefault();
});
}
});
</script>

<style>
    .item {
    padding: 10px;
    margin: 5px;
    background-color: #f0f0f0;
    cursor: move;
}
</style>

        {#each items as item, index (item.id)}
        {#key item.id}
        <div
            class="item"
            draggable="true"
            on:dragstart={handleDragStart}
            on:drop={handleDrop}
            on:dragover={(event) => event.preventDefault()}
            data-index={index}
            id={item.id}
            transition:fade
        >
            {item.text}
        </div>
        {/key}
{/each}
