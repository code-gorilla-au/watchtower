import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { TimeSince } from "$lib/hooks/time.svelte";
import { SvelteDate } from "svelte/reactivity";

describe("Hooks: time", () => {
	beforeEach(() => {
		vi.useFakeTimers();
	});
	afterEach(() => {
		vi.resetAllMocks();
	});

	it("should return reactive time since", () => {
		const tt = new TimeSince(new SvelteDate());
		expect(tt.date).toBe("less than a minute");
	});
});
