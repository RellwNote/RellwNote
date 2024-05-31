<script lang="ts">
    import ContentListComponent from "../components/ContentListComponent.svelte";
    import Router from "svelte-spa-router";
    import routes from "../router/routes";
</script>

<div class="container">
    <div class="content">
        <div class="contentLogo">RellwNoteX</div>
        <ContentListComponent />
    </div>
    <div class="mainView">
        <Router {routes} />
    </div>
</div>

<style>
    .container {
        height: 100vh;

        display: flex;

    }
    .content {
        flex: 1;

        display: flex;
        flex-direction: column;
        align-items: center;
        background-color: #295560;
    }
    .mainView {
        flex: 3;
        background: #fff7f7;
    }
    .contentLogo {
        color: white;
        font-size: 2rem;
        height: 5rem;
        line-height: 5rem;
    }
</style>