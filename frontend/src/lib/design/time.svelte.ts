import type { SvelteDate } from "svelte/reactivity";
import { formatDate } from "$design/formats";

export class TimeSince {
	readonly #date: SvelteDate;

	readonly #poll: number;
	readonly date: string;

	constructor(date: SvelteDate) {
		this.#date = date;
		this.date = $derived(formatDate(this.#date));

		this.#poll = setInterval(this.start, 1000);
	}

	private start() {
		this.#date.setTime(Date.now());
	}

	stop() {
		clearInterval(this.#poll);
	}

	setDate(date: SvelteDate) {
		this.#date.setDate(date.getDate());
	}
}
