<script lang="ts">
    export let color: string = "info";
    export let size: string = "medium";

    let btnRef: any;


</script>

<button class="WBtn {color} {size}" bind:this={btnRef}><slot></slot></button>

<style>
    .WBtn {
        border-radius: 0.4rem;

        font-size: 1.1rem;
        color: #efeaea;
    }

    .WBtn:hover {
        opacity: 0.8;
    }

    .info {
        background-color: #7c7cd1;
    }

    .success {
        background-color: #4000ff;
    }

    .warning {
         background-color: #98ec53;
    }

    .error {
          background-color: #FF0000;
    }

    .small {
        width: 6rem;
        height: 2.5rem;
    }

    .large {
        width: 15rem;
        height: 5rem;

        font-size: 1.8rem;
    }

    .medium {
        width: 10rem;
        height: 3rem;
        font-size: 1.3rem;
    }
</style>