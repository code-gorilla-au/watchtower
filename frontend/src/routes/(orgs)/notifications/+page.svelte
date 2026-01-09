<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { EmptySlate } from "$components/empty_slate";

	const { data } = $props();

	let notifications = $derived(data.notifications);
</script>

<div class="page-container">
	<PageTitle title="Notifications" subtitle="Unread notifications across all orgs" />

	{#if notifications.length === 0}
		<EmptySlate title="No new notifications"></EmptySlate>
	{:else}
		{#each notifications as notification (notification.id)}
			<div>
				<p>{notification.content}</p>
				<p>{notification.created_at}</p>
			</div>
		{/each}
	{/if}
</div>
