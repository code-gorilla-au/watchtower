<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { EmptySlate } from "$components/empty_slate";
	import * as Card from "$lib/components/ui/card";
	import { Badge } from "$lib/components/ui/badge";
	import { Button } from "$lib/components/ui/button";
	import { Check, Bell, Inbox } from "@lucide/svelte";
	import { formatDate, toSentenceCase } from "$lib/hooks/formats";
	import { notificationSvc } from "$lib/watchtower";
	import { invalidateAll } from "$app/navigation";

	const { data } = $props();

	let notifications = $derived(data.notifications);

	async function markAllAsRead() {
		await notificationSvc.markAllAsRead();
		await invalidateAll();
	}

	async function markAsRead(id: number) {
		await notificationSvc.markAsRead(id);
		await invalidateAll();
	}
</script>

<div class="page-container flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<PageTitle title="Notifications" subtitle="Unread notifications across all orgs" />
		{#if notifications.length > 0}
			<Button variant="outline" size="sm" onclick={markAllAsRead}>
				<Check class="mr-2 h-4 w-4" />
				Mark all read
			</Button>
		{/if}
	</div>

	{#if notifications.length === 0}
		<EmptySlate
			title="No new notifications"
			description="You're all caught up! Check back later for updates."
		>
			<div class="mt-4">
				<Inbox class="h-12 w-12 text-muted-foreground/20" />
			</div>
		</EmptySlate>
	{:else}
		<div class="flex flex-col gap-2">
			{#each notifications as notification (notification.id)}
				<Card.Root class="overflow-hidden p-0">
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
			{/each}
		</div>
	{/if}
</div>
