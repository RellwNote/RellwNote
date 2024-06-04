<script lang="ts">
    import WBtn from "../components/WBtn.svelte";
    import WInput from "../components/WInput.svelte";
</script>
<div class="indexWindow">
    <div class="serverStart">
        <WInput type="text" placeholder=":8888"/>
        <WBtn color="info">启动</WBtn>
    </div>
    <div class="noteDetail">

    </div>
</div>


<style>
    .indexWindow {
        display: flex;
        flex-direction: column;
        height: 100vh;
    }

    .serverStart {
        flex: 1;

        display: flex;
        align-items: center;
        justify-content: space-around;
    }

    .noteDetail {
        flex: 5;
        background-color: #5cabbb;
    }
</style>