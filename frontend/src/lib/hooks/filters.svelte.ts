export type FilterFunc<T> = (item: T) => boolean;

export interface FilterState<T> {
	initData: T[];
	filteredData: T[];
}

export class SimpleFilter<T> {
	#internal: FilterState<T>;
	#filterFn?: FilterFunc<T> = undefined;

	readonly data: T[];

	constructor(data: T[], filterFn?: FilterFunc<T>) {
		this.#internal = $state({ filteredData: [], initData: [...data] });
		this.#filterFn = filterFn;
		this.data = $derived(this.#internal.filteredData);

		this.applyFilter();

		$effect(() => {
			this.applyFilter();
		});
	}

	/**
	 * Clears the current filter function and resets the filter state.
	 *
	 */
	clear() {
		this.#filterFn = undefined;
		this.applyFilter();
	}

	/**
	 * Applies the provided filter function to the current dataset and filters it based on the criteria defined in the function.
	 */
	filterBy(filterFn: FilterFunc<T>) {
		this.#filterFn = filterFn;
		this.applyFilter();
	}

	private applyFilter() {
		if (!this.#filterFn) {
			this.resetInitState();
			return;
		}

		this.#internal.filteredData.splice(
			0,
			this.#internal.filteredData.length,
			...this.#internal.initData.filter(this.#filterFn)
		);
	}

	private resetInitState() {
		this.#internal.filteredData.splice(
			0,
			this.#internal.filteredData.length,
			...this.#internal.initData
		);
	}
}
