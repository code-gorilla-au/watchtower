import { describe, expect, it } from "vitest";
import { SimpleFilter, TagsFilter } from "$lib/hooks/filters.svelte";

describe("filters", () => {
	describe("SimpleFilter()", () => {
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

			expect(filter.data).toEqual([2, 3]);
		});
	});

	describe("TagsFilter()", () => {
		const validTags = [
			{ name: "hello", value: "world" },
			{ name: "goodbye", value: "planet" },
			{ name: "baz", value: "world" }
		];

		it("should have a list of tags available", () => {
			const filter = new TagsFilter(validTags, "value");
			expect(filter.tags).toContain("world");
			expect(filter.tags).toContain("planet");
		});

		it("should return init data without filter", () => {
			const filter = new TagsFilter(validTags, "value");
			expect(filter.data).toEqual(validTags);
		});

		it("should return filtered data with filter", () => {
			const filter = new TagsFilter(validTags, "value");
			expect(filter.data).toEqual(validTags);

			filter.filterByTag("planet");
			expect(filter.data).toEqual([{ name: "goodbye", value: "planet" }]);
		});

		it("should have current tag applied when filter is applied", () => {
			const filter = new TagsFilter(validTags, "value");

			filter.filterByTag("planet");
			expect(filter.currentTag).toEqual("planet");
		});

		it("should change filter after applying twice", () => {
			const filter = new TagsFilter(validTags, "value");

			filter.filterByTag("planet");
			expect(filter.data).toEqual([{ name: "goodbye", value: "planet" }]);

			filter.filterByTag("world");
			expect(filter.data).toEqual([
				{ name: "hello", value: "world" },
				{ name: "baz", value: "world" }
			]);
		});

		it("should filter based on multiple tags", () => {
			const multiTags = [
				{
					name: "hello",
					value: "world",
					tag: "tag1"
				},
				{
					name: "bin",
					value: "baz",
					tag: "tag2"
				},
				{
					name: "hawk",
					value: "flash",
					tag: "tag1"
				}
			];

			const filter = new TagsFilter(multiTags, ["value", "tag"]);

			filter.filterByTag("tag2");
			expect(filter.data).toEqual([
				{
					name: "bin",
					value: "baz",
					tag: "tag2"
				}
			]);
		});
	});
});
