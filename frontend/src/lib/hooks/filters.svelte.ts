export type FilterFunc<T> = (item: T) => boolean;

export interface FilterState<T> {
	data: T[];
	filterFn?: FilterFunc<T>;
}

export class SimpleFilter<T> {
	#internal: FilterState<T>;

	readonly data: T[];

	constructor(data: T[], filterFn?: FilterFunc<T>) {
		this.#internal = $state({ data, filterFn });

		this.data = $derived(this.#internal.data);
	}
}
