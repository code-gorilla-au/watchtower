import { describe, expect, it } from "vitest";
import { SimpleFilter } from "$lib/hooks/filters.svelte";

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
	});
});
