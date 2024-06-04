<script lang="ts">
    export let type: string;
    export let placeholder: string;
</script>

<input type={type} placeholder={placeholder} class="WInput">

<style>
    .WInput {
        width: 25rem;
        height: 2.5rem;
        font-size: 1.2rem;

        border-radius: 0.5rem;
    }

    .WInput::placeholder {
        font-style: italic;
    }
</style>