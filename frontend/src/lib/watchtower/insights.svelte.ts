import { insights } from "$lib/wailsjs/go/models";
import { STALE_TIMEOUT_MINUTES } from "$lib/watchtower/types";
import {
	GetPullRequestInsightsByOrg,
	GetSecurityInsightsByOrg
} from "$lib/wailsjs/go/watchtower/Service";

export class InsightsService {
	#defaultTimeWindow = "90";
	#sec?: insights.SecurityInsights;
	#pr?: insights.PullRequestInsights;
	#lastSync?: number;

	hasInsights: boolean;

	constructor() {
		this.#sec = $state(undefined);
		this.#pr = $state(undefined);
		this.hasInsights = $derived((!!this.#sec || !!this.#pr) ?? false);
	}

	get window() {
		return this.#defaultTimeWindow;
	}

	get prInsights() {
		return this.#pr;
	}

	get secInsights() {
		return this.#sec;
	}

	async refresh(orgId: number) {
		await this.forceGetInsights(orgId, this.#defaultTimeWindow);

		return { sec: this.#sec, pr: this.#pr };
	}

	async getInsights(orgId: number) {
		if (!this.isStale()) {
			return { sec: this.#sec, pr: this.#pr };
		}

		return this.refresh(orgId);
	}

	private async forceGetInsights(orgId: number, timeWindow: string) {
		await this.forceGetSec(orgId, timeWindow);
		await this.forceGetPR(orgId, timeWindow);
	}

	private async forceGetSec(orgId: number, timeWindow: string) {
		this.#sec = await GetSecurityInsightsByOrg(orgId, timeWindow);
		this.#lastSync = Date.now();
	}

	private async forceGetPR(orgId: number, timeWindow: string) {
		this.#pr = await GetPullRequestInsightsByOrg(orgId, timeWindow);

		this.#lastSync = Date.now();
	}

	private isStale() {
		if (!this.#lastSync) {
			return true;
		}

		const diff = (Date.now() - this.#lastSync) / (1000 * 60);
		return diff > STALE_TIMEOUT_MINUTES;
	}
}
