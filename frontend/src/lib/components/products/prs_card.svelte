<script lang="ts">
	import { Card, CardAction, CardHeader, CardTitle } from "$components/ui/card/index.js";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { Badge } from "$components/ui/badge/index.js";
	import { formatDate } from "$design/formats";
	import { SquareArrowOutUpRight } from "@lucide/svelte";
	import { Button } from "$components/ui/button";

	type Props = {
		pr: watchtower.PullRequestDTO;
	};

	let { pr }: Props = $props();
</script>

<Card>
	<CardHeader class="flex items-center justify-between">
		<CardTitle>
			<span class="capitalize">{pr.title}</span>
		</CardTitle>
		<CardAction>
			<Button href={pr.url} target="_blank" size="icon" variant="ghost">
				<SquareArrowOutUpRight />
			</Button>
		</CardAction>
	</CardHeader>

	<div class="px-3">
		<div class="mb-2 flex items-center justify-between">
			<p class="text-sm text-muted-foreground">Last updated</p>
			<p>{formatDate(pr.updated_at)}</p>
		</div>
		<div>
			<Badge>{pr.repository_name}</Badge>
			<Badge>{pr.state}</Badge>
		</div>
	</div>
</Card>
