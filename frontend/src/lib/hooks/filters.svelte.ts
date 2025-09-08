import { SvelteSet } from "svelte/reactivity";

export type FilterFunc<T> = (item: T) => boolean;

/**
 * Represents the state of a filter operation, containing both the original
 * unfiltered dataset and the filtered results.
 */
export interface FilterState<T> {
	/**
	 * Original data stored. Used to apply filter.
	 */
	initData: T[];
	/**
	 * An array that holds the filtered dataset based on specific criteria.
	 */
	filteredData: T[];
}

export type FilterTag<T> = keyof T;
export type FilterTagValue<T> = T[keyof T];

/**
 * A utility class to filter a dataset based on a given filtering function.
 * This class maintains an internal filter state and provides methods to apply or clear filters.
 */
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

type TagField<T> = FilterTag<T> | FilterTag<T>[];

/**
 * Represents a filter that can be used to filter data by specific tags.
 */
export class TagsFilter<T extends object> {
	#filter: SimpleFilter<T>;
	#currentTag?: FilterTagValue<T>;
	readonly #tagField: TagField<T>;
	readonly tags: FilterTagValue<T>[];
	readonly currentTag: FilterTagValue<T> | undefined;
	readonly data: T[];

	constructor(data: T[], tagField: TagField<T>) {
		this.#currentTag = $state(undefined);
		this.#tagField = tagField;
		this.tags = $state(this.generateTags(tagField, data));
		this.#filter = new SimpleFilter(data, this.filterFn.bind(this));

		this.currentTag = $derived(this.#currentTag);
		this.data = $derived(this.#filter.data);
	}

	filterByTag(tag: FilterTagValue<T>) {
		this.#currentTag = tag;
		this.#filter.filterBy(this.filterFn.bind(this));
	}

	reset() {
		this.#currentTag = undefined;
		this.#filter.clear();
	}

	private generateTags(tags: TagField<T>, data: T[]) {
		const set = new SvelteSet<FilterTagValue<T>>();

		for (const item of data) {
			if (!Array.isArray(tags)) {
				const value = item[tags];
				set.add(value);
				continue;
			}

			for (const tag of tags) {
				const value = item[tag];
				set.add(value);
			}
		}

		return Array.from(set);
	}

	private filterFn(item: T) {
		if (!this.#currentTag) {
			return true;
		}

		if (!Array.isArray(this.#tagField)) {
			return item[this.#tagField] === this.#currentTag;
		}

		for (const tag of this.#tagField) {
			if (item[tag] === this.#currentTag) {
				return true;
			}
		}

		return false;
	}
}
