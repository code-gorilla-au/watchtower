import { SvelteDate } from "svelte/reactivity";
import { formatDistanceToNow } from "date-fns";

/**
 * The TimeSince class calculates and tracks the time elapsed since a given date.
 * It provides functionality to start and stop periodic updates of the elapsed time.
 */
export class TimeSince {
	readonly #date: SvelteDate;
	#timeSince: string;

	#poll?: NodeJS.Timeout;
	readonly date: string;

	constructor(date: SvelteDate) {
		this.#date = date;
		this.#timeSince = $state(this.differenceInTime());
		this.date = $derived(this.#timeSince.toString());
	}

	/**
	 * Starts the interval for executing the `updateTimeSince` method every second.
	 * This method sets up a poll using `setInterval` to repeatedly call the `updateTimeSince` method.
	 */
	start() {
		this.#poll = setInterval(this.updateTimeSince.bind(this), 1000);
	}

	/**
	 * Stops the polling process by clearing the interval.
	 */
	stop() {
		clearInterval(this.#poll);
	}

	/**
	 * Sets the date to the provided SvelteDate instance.
	 */
	setDate(time: SvelteDate) {
		this.#date.setTime(time.getTime());
	}

	private updateTimeSince() {
		this.#timeSince = this.differenceInTime();
	}

	private differenceInTime() {
		return formatDistanceToNow(this.#date);
	}
}
