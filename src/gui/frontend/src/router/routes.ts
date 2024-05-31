import GitPushView from "../views/GitPushView.svelte";
import SettingView from "../views/SettingView.svelte";
import TOCConfigView from "../views/TOCConfigView.svelte";
import IndexView from "../views/IndexView.svelte";


const routes = {
    '/': IndexView,
    '/tocconfigview': TOCConfigView,
    '/settingview': SettingView,
    '/gitpushview': GitPushView,
}

export default routes;