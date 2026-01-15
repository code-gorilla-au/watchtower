import { describe, it, expect, vi, beforeEach } from "vitest";
import * as wails from "$lib/wailsjs/go/watchtower/Service";
import { insights } from "$lib/wailsjs/go/models";
import { InsightsService } from "$lib/watchtower/insights.svelte";

describe("InsightsService", () => {
	const spyGetPR = vi.spyOn(wails, "GetPullRequestInsightsByOrg");
	const spyGetSec = vi.spyOn(wails, "GetSecurityInsightsByOrg");

	const mockPR = new insights.PullRequestInsights({
		merged: 10,
		closed: 2,
		open: 5,
		avgDaysToMerge: 3.5,
		minDaysToMerge: 1,
		maxDaysToMerge: 10
	});

	const mockSec = new insights.SecurityInsights({
		fixed: 8,
		open: 3,
		avgDaysToFix: 5.2,
		minDaysToFix: 1,
		maxDaysToFix: 15
	});

	beforeEach(() => {
		vi.clearAllMocks();
		vi.useFakeTimers();
	});

	describe("initial state", () => {
		it("should have default time window", () => {
			const service = new InsightsService();
			expect(service.window).toBe("90");
		});

		it("should not have insights initially", () => {
			const service = new InsightsService();
			expect(service.hasInsights).toBe(false);
			expect(service.prInsights).toBeUndefined();
			expect(service.secInsights).toBeUndefined();
		});
	});

	describe("getInsights()", () => {
		beforeEach(() => {
			spyGetPR.mockResolvedValue(mockPR);
			spyGetSec.mockResolvedValue(mockSec);
		});

		it("should fetch insights if not already present", async () => {
			const service = new InsightsService();
			const result = await service.getInsights(1);

			expect(spyGetPR).toHaveBeenCalledWith(1, "90");
			expect(spyGetSec).toHaveBeenCalledWith(1, "90");
			expect(result.pr).toEqual(mockPR);
			expect(result.sec).toEqual(mockSec);
			expect(service.hasInsights).toBe(true);
		});

		it("should not refetch if not stale", async () => {
			const service = new InsightsService();
			await service.getInsights(1);
			expect(spyGetPR).toHaveBeenCalledTimes(1);

			await service.getInsights(1);
			expect(spyGetPR).toHaveBeenCalledTimes(1);
		});

		it("should refetch if stale", async () => {
			const service = new InsightsService();
			await service.getInsights(1);
			expect(spyGetPR).toHaveBeenCalledTimes(1);

			vi.setSystemTime(new Date(Date.now() + 3 * 60 * 1000));

			await service.getInsights(1);
			expect(spyGetPR).toHaveBeenCalledTimes(2);
		});
	});

	describe("refresh()", () => {
		beforeEach(() => {
			spyGetPR.mockResolvedValue(mockPR);
			spyGetSec.mockResolvedValue(mockSec);
		});

		it("should force refetch even if not stale", async () => {
			const service = new InsightsService();
			await service.getInsights(1);
			expect(spyGetPR).toHaveBeenCalledTimes(1);

			await service.refresh(1);
			expect(spyGetPR).toHaveBeenCalledTimes(2);
		});
	});

	describe("hasInsights", () => {
		it("should be true if only PR insights are present", async () => {
			spyGetPR.mockResolvedValue(mockPR);

			const service = new InsightsService();
			await service.getInsights(1);

			expect(service.hasInsights).toBe(true);
			expect(service.prInsights).toBeDefined();
		});

		it("should be true if only Sec insights are present", async () => {
			spyGetSec.mockResolvedValue(mockSec);

			const service = new InsightsService();
			await service.getInsights(1);

			expect(service.hasInsights).toBe(true);
			expect(service.secInsights).toBeDefined();
		});
	});
});
