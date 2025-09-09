<script lang="ts">
	import { Card, CardAction, CardHeader, CardTitle } from "$components/ui/card/index.js";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { Badge } from "$components/ui/badge/index.js";
	import { formatDate, truncate } from "$lib/hooks/formats";
	import { SquareArrowOutUpRight } from "@lucide/svelte";
	import { Button } from "$components/ui/button";
	import { OpenExternalURL } from "$lib/wailsjs/go/main/App";

	type Props = {
		pr: watchtower.PullRequestDTO;
	};

	let { pr }: Props = $props();

	async function routeToPr(e: Event) {
		e.preventDefault();
		await OpenExternalURL(pr.url);
	}
</script>

<Card>
	<CardHeader class="flex items-center justify-between">
		<CardTitle>
			<span class="capitalize">{truncate(pr.title)}</span>
		</CardTitle>
		<CardAction>
			<Button onclick={routeToPr} size="icon" variant="ghost">
				<SquareArrowOutUpRight />
			</Button>
		</CardAction>
	</CardHeader>

	<div class="px-3">
		<div class="card-row">
			<p class="row-label">Author</p>
			<p>{pr.author}</p>
		</div>
		<div class="card-row">
			<p class="row-label">Status</p>
			<p class="flex-1 lowercase">{pr.state}</p>
		</div>
		<div class="card-row">
			<p class="row-label">Created</p>
			<p>{formatDate(pr.created_at)}</p>
		</div>
		<div class="card-row">
			<p class="row-label">Repository</p>
			<p class="flex-1 lowercase">{pr.repository_name}</p>
		</div>
		<div>
			<Badge>{pr.tag}</Badge>
		</div>
	</div>
</Card>

<style lang="postcss">
	@reference "$design";

	.card-row {
		@apply mb-2 flex items-baseline;
	}

	.row-label {
		@apply w-1/2 text-xs text-muted-foreground;
	}
</style>
