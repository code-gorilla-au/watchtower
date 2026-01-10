<script lang="ts">
	import { fade, fly } from "svelte/transition";
	import { flip } from "svelte/animate";
	import { PageTitle } from "$components/page_title";
	import { EmptySlate } from "$components/empty_slate";
	import * as Card from "$lib/components/ui/card";
	import { Badge } from "$lib/components/ui/badge";
	import { Button } from "$lib/components/ui/button";
	import { Check, Bell, Inbox, Search } from "@lucide/svelte";
	import { formatDate, toSentenceCase } from "$lib/hooks/formats";
	import { notificationSvc, orgSvc } from "$lib/watchtower";
	import { goto, invalidateAll } from "$app/navigation";
	import { resolve } from "$app/paths";
	import { SimpleFilter } from "$lib/hooks/filters.svelte";
	import { BaseInput } from "$components/base_input";
	import { notifications } from "$lib/wailsjs/go/models";

	const { data } = $props();

	let unreadNotifications = $derived(data.notifications);

	let searchState = $state("");

	function applySearchFilter(notification: notifications.Notification) {
		return notification.content.toLowerCase().includes(searchState.toLowerCase());
	}

	const searchFilter = $derived(new SimpleFilter(unreadNotifications, applySearchFilter));

	async function markAllAsRead() {
		await notificationSvc.markAllAsRead();
		await invalidateAll();
	}

	async function markAsRead(id: number) {
		await notificationSvc.markAsRead(id);
		await invalidateAll();
	}

	async function routeToOrgDashboard(orgId: number) {
		await orgSvc.setDefault(orgId);
		await goto(resolve("/dashboard"));
	}
</script>

<div class="page-container">
	<div class="flex items-center justify-between">
		<PageTitle title="Notifications" subtitle="Unread notifications across all orgs" />
		{#if searchFilter.data.length > 0}
			<Button variant="outline" size="sm" onclick={markAllAsRead}>
				<Check class="mr-2 h-4 w-4" />
				Mark all read
			</Button>
		{/if}
	</div>

	{#if unreadNotifications.length === 0}
		<EmptySlate
			title="No new notifications"
			description="You're all caught up! Check back later for updates."
		>
			<div class="mt-4">
				<Inbox class="h-12 w-12 text-muted-foreground/20" />
			</div>
		</EmptySlate>
	{:else}
		<div class="mb-2 flex w-full justify-end">
			<div class="flex items-center gap-2">
				<Search class="" />
				<BaseInput bind:value={searchState} placeholder="Filter notifications" />
			</div>
		</div>
		<div class="flex flex-col gap-2">
			{#each searchFilter.data as notification (notification.id)}
				<div animate:flip in:fade={{ duration: 200 }} out:fly={{ x: 100, duration: 200 }}>
					<Card.Root
						onclick={(e) => {
							e.preventDefault();

							routeToOrgDashboard(notification.organisation_id);
						}}
						class="overflow-hidden p-0 hover:cursor-pointer"
					>
						<div class="flex items-center gap-4 p-4">
							<div class="shrink-0">
								<div
									class="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10"
								>
									<Bell class="h-5 w-5 text-primary" />
								</div>
							</div>
							<div class="min-w-0 flex-1">
								<div class="mb-1 flex items-center gap-2">
									<Badge variant="secondary" class="text-xs">
										{toSentenceCase(notification.type)}
									</Badge>
									<span class="text-xs text-muted-foreground">
										{formatDate(notification.created_at)}
									</span>
								</div>
								<p class="truncate text-sm font-medium">
									{notification.content}
								</p>
							</div>
							<div class="shrink-0">
								<Button
									variant="ghost"
									size="icon"
									class="h-8 w-8"
									onclick={(e) => {
										e.stopImmediatePropagation();

										markAsRead(notification.id);
									}}
									title="Mark as read"
								>
									<Check class="h-4 w-4" />
								</Button>
							</div>
						</div>
					</Card.Root>
				</div>
			{/each}
		</div>
	{/if}
</div>
