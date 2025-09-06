<script lang="ts">
	import { Badge } from "$components/ui/badge";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { TagsFilter } from "$lib/hooks/filters.svelte";

	type Props = {
		prs: watchtower.PullRequestDTO[];
	};

	let { prs }: Props = $props();

	const prsFilter = new TagsFilter(prs, "repository_name");

	function styleCurrentTag(currentTag: string, tag: string) {
		if (currentTag === tag) {
			return "default";
		}

		return "outline";
	}
</script>

<div class="">
	<h2 class="mb-1 text-xl text-muted-foreground">Pull Requests ({prs.length})</h2>
	<div class="flex gap-2">
		{#each prsFilter.tags as tag (tag)}
			<button
				onclick={(e: Event) => {
					e.stopPropagation();
					prsFilter.filterByTag(tag);
				}}
			>
				<Badge variant={styleCurrentTag(prsFilter.currentTag, tag)} class="text-xs">
					{tag}
				</Badge>
			</button>
		{/each}
		{#if prsFilter.currentTag}
			<button
				class=""
				onclick={(e: Event) => {
					e.stopPropagation();
					prsFilter.reset();
				}}
			>
				clear
			</button>
		{/if}
	</div>
</div>
