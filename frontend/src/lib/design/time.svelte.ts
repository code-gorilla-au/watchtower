import { SvelteDate } from "svelte/reactivity";
import { differenceInSeconds } from "date-fns";

export class TimeSince {
	readonly #date: SvelteDate;

	readonly #poll: number;
	readonly date: string;

	constructor(date: SvelteDate) {
		this.#date = date;
		this.date = $derived(differenceInSeconds(this.#date, new SvelteDate()).toString());

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
