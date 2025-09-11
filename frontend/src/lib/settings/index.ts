import { Settings } from "./env.svelte";
import { ConfigFileLocation, Version } from "$lib/wailsjs/go/main/App";

const appVersion = await Version();
const appConfigDir = await ConfigFileLocation();

export const settingsSvc = new Settings(appVersion, appConfigDir);
