<script lang="ts">
    import WBtn from "../components/WBtn.svelte";
    import TOCList from "../components/TOCList.svelte";
</script>

<div class="TOCWindow">
    <div class="TOCContent">
        <TOCList />
    </div>
    <div class="TOCOperate">
        <WBtn color="info" size="medium">重设</WBtn>
        <WBtn color="success" size="medium">保存</WBtn>
    </div>
</div>

<style>
    .TOCWindow {
        height: 100vh;

        display: flex;
        flex-direction: column;
    }

    .TOCContent {
        flex: 5;
        background-color: #5cabbb;
    }

    .TOCOperate{
        flex: 1;

        display: flex;
        align-items: center;
        justify-content: space-around;
    }

</style>