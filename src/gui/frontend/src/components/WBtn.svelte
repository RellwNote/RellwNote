<script lang="ts">
    import {onMount} from "svelte";
    export let color: string = "info";


    onMount(() => {
        let btnDom = document.getElementById("WBtn");
        switch (color) {
            case "info":
                btnDom.style.backgroundColor = "#7c7cd1";
                break;
            case "success":
                btnDom.style.backgroundColor = "#FF0000";
                break;
        }

    })

</script>

<button class="btn WBtn" id="WBtn"><slot></slot></button>

<style>
    .WBtn {
        width: 7rem;
        height: 2.7rem;

        border-radius: 0.4rem;

        font-size: 1.1rem;
        color: #efeaea;
    }

    .WBtn:hover {
        opacity: 0.8;
    }
</style>