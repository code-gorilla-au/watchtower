import { describe, expect, it } from "vitest";
import { SimpleFilter } from "$lib/hooks/filters.svelte";
import { flushSync } from "svelte";

describe("filters", () => {
	describe("simple filters", () => {
		it("should return an empty data array", () => {
			const filter = new SimpleFilter([]);

			expect(filter.data).toEqual([]);
		});

		it("should return an array with data", () => {
			const filter = new SimpleFilter([1, 2, 3]);
			expect(filter.data).toEqual([1, 2, 3]);
		});

		it("it should apply filter if provided", () => {
			const filter = new SimpleFilter([1, 2, 3], (item) => {
				return item > 1;
			});

			expect(filter.data).toEqual([2, 3]);
		});

		it("clearing the filter should return original data", () => {
			const filter = new SimpleFilter([1, 2, 3], (item) => {
				return item > 1;
			});
			expect(filter.data).toEqual([2, 3]);

			filter.clear();
			expect(filter.data).toEqual([1, 2, 3]);
		});

		it("it should apply filter if provided", () => {
			const filter = new SimpleFilter([1, 2, 3]);
			expect(filter.data).toEqual([1, 2, 3]);

			filter.filterBy((item) => {
				return item > 1;
			});
			flushSync();

			expect(filter.data).toEqual([2, 3]);
		});
	});
});
