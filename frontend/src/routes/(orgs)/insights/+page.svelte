<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";
	import * as Card from "$components/ui/card";
	import {
		GitPullRequest,
		ShieldCheck,
		Clock,
		CheckCircle2,
		XCircle,
		AlertCircle,
		RefreshCcw
	} from "@lucide/svelte";
	import { invalidateAll } from "$app/navigation";
	import { Button } from "$components/ui/button";

	let { data }: PageProps = $props();
	let org = $derived(data.organisation);
	let prInsights = $derived(data.insights.pr);
	let secInsights = $derived(data.insights.sec);
	let insightWindow = $derived(data.insights.window);

	let refreshing = $state(false);

	async function handleRefresh() {
		refreshing = true;
		await invalidateAll();
		refreshing = false;
	}

	function formatDays(days: number) {
		if (days === undefined || days === null) return "N/A";
		return days.toFixed(1) + " days";
	}
</script>

<div class="page-container">
	<PageTitle title="Insights" subtitle={org?.description || org?.friendly_name}>
		<Button variant="outline" size="sm" onclick={handleRefresh} disabled={refreshing}>
			<RefreshCcw class="mr-2 h-4 w-4 {refreshing ? 'animate-spin' : ''}" />
			Refresh
		</Button>
	</PageTitle>

	<div class="mt-6 grid grid-cols-1 gap-6 md:grid-cols-2">
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<GitPullRequest class="h-5 w-5" />
					Pull Request Insights
				</Card.Title>
				<Card.Description>
					Statistics for merged, closed, and open pull requests over the last {insightWindow}
					days.
				</Card.Description>
			</Card.Header>
			<Card.Content>
				{#if prInsights}
					<div class="grid grid-cols-2 gap-4">
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Total Merged</p>
							<div class="flex items-center gap-2">
								<CheckCircle2 class="h-4 w-4 text-green-500" />
								<p class="text-2xl font-bold">{prInsights.merged}</p>
							</div>
						</div>
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Total Closed</p>
							<div class="flex items-center gap-2">
								<XCircle class="h-4 w-4 text-destructive" />
								<p class="text-2xl font-bold">{prInsights.closed}</p>
							</div>
						</div>
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Currently Open</p>
							<div class="flex items-center gap-2">
								<AlertCircle class="h-4 w-4 text-accent" />
								<p class="text-2xl font-bold">{prInsights.open}</p>
							</div>
						</div>
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Avg. Time to Merge</p>
							<div class="flex items-center gap-2">
								<Clock class="h-4 w-4 text-muted-foreground" />
								<p class="text-2xl font-bold">
									{formatDays(prInsights.avgDaysToMerge)}
								</p>
							</div>
						</div>
					</div>

					<div class="mt-6 grid grid-cols-2 gap-4 border-t pt-6">
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Min. Time to Merge</p>
							<p class="text-lg font-semibold">
								{formatDays(prInsights.minDaysToMerge)}
							</p>
						</div>
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Max. Time to Merge</p>
							<p class="text-lg font-semibold">
								{formatDays(prInsights.maxDaysToMerge)}
							</p>
						</div>
					</div>
				{:else}
					<div class="flex justify-center p-8">
						<p class="text-muted-foreground">No PR insights available</p>
					</div>
				{/if}
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<ShieldCheck class="h-5 w-5" />
					Security Insights
				</Card.Title>
				<Card.Description>
					Vulnerability status and remediation metrics for the last {insightWindow} days.
				</Card.Description>
			</Card.Header>
			<Card.Content>
				{#if secInsights}
					<div class="grid grid-cols-2 gap-4">
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Fixed Vulnerabilities</p>
							<div class="flex items-center gap-2">
								<CheckCircle2 class="h-4 w-4 text-green-500" />
								<p class="text-2xl font-bold">{secInsights.fixed}</p>
							</div>
						</div>
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Open Vulnerabilities</p>
							<div class="flex items-center gap-2">
								<AlertCircle class="h-4 w-4 text-destructive" />
								<p class="text-2xl font-bold">{secInsights.open}</p>
							</div>
						</div>
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Avg. Time to Fix</p>
							<div class="flex items-center gap-2">
								<Clock class="h-4 w-4 text-muted-foreground" />
								<p class="text-2xl font-bold">
									{formatDays(secInsights.avgDaysToFix)}
								</p>
							</div>
						</div>
					</div>

					<div class="mt-6 grid grid-cols-2 gap-4 border-t pt-6">
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Min. Time to Fix</p>
							<p class="text-lg font-semibold">
								{formatDays(secInsights.minDaysToFix)}
							</p>
						</div>
						<div class="space-y-1">
							<p class="text-sm text-muted-foreground">Max. Time to Fix</p>
							<p class="text-lg font-semibold">
								{formatDays(secInsights.maxDaysToFix)}
							</p>
						</div>
					</div>
				{:else}
					<div class="flex justify-center p-8">
						<p class="text-muted-foreground">No security insights available</p>
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
	</div>
</div>
