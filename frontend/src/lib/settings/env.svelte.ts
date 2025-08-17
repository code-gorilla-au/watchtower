import { browser } from "$app/environment";

const SETTINGS_KEY = "--app-settings";

type AppSettings = {
	theme: "light" | "dark";
};

export class Settings {
	#internal: AppSettings = $state({
		theme: "dark"
	});

	readonly theme = $derived(this.#internal.theme);

	constructor() {
		this.init();
		this.applyTheme(this.#internal.theme);
	}

	initTheme() {
		this.applyTheme(this.#internal.theme);
	}

	setTheme(theme: AppSettings["theme"]) {
		this.#internal.theme = theme;
		this.save();
		this.applyTheme(theme);
	}

	private applyTheme(theme: AppSettings["theme"]) {
		if (!browser) {
			return;
		}

		if (theme === "dark") {
			document.documentElement.classList.add(theme);
			return;
		}

		document.documentElement.classList.remove(theme);
	}

	private init() {
		const settings = localStorage.getItem(SETTINGS_KEY);
		if (!settings) {
			return;
		}

		const appSettings = JSON.parse(settings) as AppSettings;
		Object.assign(this.#internal, appSettings);
	}

	private save() {
		localStorage.setItem(SETTINGS_KEY, JSON.stringify(this.#internal));
	}
}
