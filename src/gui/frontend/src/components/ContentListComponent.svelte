<script lang="ts">
    import {push} from "svelte-spa-router";

    const checkMenuToIndexView = () => {push('/')}
    const checkMenuToTOCConfigView = () => {push('/tocconfigview')}
    const checkMenuToGitPushView = () => {push('/gitpushview')}
    const checkMenuToSettingView = () => {push('/settingview')}

</script>

<ul>
    <li on:click={checkMenuToIndexView}>首页</li>
    <li on:click={checkMenuToTOCConfigView}>目录配置</li>
    <li on:click={checkMenuToGitPushView}>Git推送</li>
    <li on:click={checkMenuToSettingView}>设置</li>
</ul>

<style>
    ul {
        width: 100%;
        list-style-type: none;
    }
    li {
        width: 100%;
        height: 3rem;
        line-height: 3rem;
        font-size: 1.2rem;
        color: white;
        text-align: center;
        border-width: 1px 0;
        border-style: solid;
        border-color: #9f7878;
    }

    li:hover {
        background-color: #31616c;
    }
</style>