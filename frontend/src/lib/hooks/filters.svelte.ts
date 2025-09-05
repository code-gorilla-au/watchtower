export type FilterFunc<T> = (item: T) => boolean;

export interface FilterState<T> {
	initData: T[];
	data: T[];
	filterFn?: FilterFunc<T>;
}

export class SimpleFilter<T> {
	#internal: FilterState<T>;

	readonly data: T[];

	constructor(data: T[], filterFn?: FilterFunc<T>) {
		this.#internal = $state({ data, initData: data, filterFn });
		this.data = $derived(this.#internal.data);

		this.applyFilter();

		$effect(() => {
			this.applyFilter();
		});
	}

	clear() {
		this.#internal.filterFn = undefined;
		this.applyFilter();
	}

	filterBy(filterFn: FilterFunc<T>) {
		this.#internal.filterFn = filterFn;
		this.applyFilter();
	}

	private applyFilter() {
		if (this.#internal.filterFn) {
			this.#internal.data = this.#internal.data.filter(this.#internal.filterFn);
			return;
		}

		this.#internal.data.splice(0, this.#internal.data.length, ...this.#internal.initData);
	}
}
