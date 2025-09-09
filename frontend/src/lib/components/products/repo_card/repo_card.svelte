<script lang="ts">
	import { Card, CardAction, CardHeader, CardTitle } from "$components/ui/card/index.js";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { Badge } from "$components/ui/badge/index.js";
	import { formatDate } from "$lib/hooks/formats";
	import { SquareArrowOutUpRight } from "@lucide/svelte";
	import { Button } from "$components/ui/button";
	import { OpenExternalURL } from "$lib/wailsjs/go/main/App";

	type Props = {
		repo: watchtower.RepositoryDTO;
	};

	let { repo }: Props = $props();

	async function routeToRepo(e: Event) {
		e.preventDefault();
		await OpenExternalURL(repo.url);
	}
</script>

<Card>
	<CardHeader class="flex items-center justify-between">
		<CardTitle>
			<span class="capitalize">{repo.name}</span>
		</CardTitle>
		<CardAction>
			<Button onclick={routeToRepo} size="icon" variant="ghost">
				<SquareArrowOutUpRight />
			</Button>
		</CardAction>
	</CardHeader>

	<div class="px-3">
		<div class="card-row">
			<p class="row-label">Last updated</p>
			<p class="flex-1">{formatDate(repo.updated_at)}</p>
		</div>
		<div>
			<Badge>{repo.owner}</Badge>
			<Badge>{repo.topic}</Badge>
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
