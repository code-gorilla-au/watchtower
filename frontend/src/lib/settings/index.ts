import { Settings } from "./env.svelte";
import { Version } from "$lib/wailsjs/go/main/App";

const appVersion = await Version();

export const settingsSvc = new Settings(appVersion);
