<script lang="ts">
	import { Card, CardTitle } from "$components/ui/card/index.js";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { Badge } from "$components/ui/badge/index.js";
	import { formatDate } from "$design/formats";

	type Props = {
		repo: watchtower.RepositoryDTO;
	};

	let { repo }: Props = $props();
</script>

<Card>
	<CardTitle class="px-3">
		<span class="capitalize">{repo.name}</span>
	</CardTitle>
	<div class="px-3">
		<div class="mb-2 flex items-center justify-between">
			<p class="text-sm text-muted-foreground">Last updated</p>
			<p>{formatDate(repo.updated_at)}</p>
		</div>
		<div>
			<Badge>{repo.owner}</Badge>
			<Badge>{repo.topic}</Badge>
		</div>
	</div>
</Card>
