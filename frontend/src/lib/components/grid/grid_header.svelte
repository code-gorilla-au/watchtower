<script lang="ts">
	import { Badge } from "$components/ui/badge";
	import { type FilterTagValue } from "$lib/hooks/filters.svelte";

	type Props = {
		title: string;
		dataLength: number;
		tags: FilterTagValue<string>[];
		currentTag?: FilterTagValue<string>;
		filterByTag: (tag: FilterTagValue<string>) => void;
		reset: () => void;
	};

	let { title, tags, currentTag, dataLength = 0, filterByTag, reset }: Props = $props();

	function styleCurrentTag(tag: FilterTagValue<string>, currentTag?: FilterTagValue<string>) {
		if (currentTag === tag) {
			return "default";
		}

		return "outline";
	}
</script>

<div class="">
	<h2 class="mb-1 text-xl text-muted-foreground">{title} ({dataLength})</h2>
	<div class="flex gap-2">
		{#each tags as tag (tag)}
			<button
				onclick={(e: Event) => {
					e.stopPropagation();
					filterByTag(tag);
				}}
			>
				<Badge variant={styleCurrentTag(tag, currentTag)} class="text-xs">
					{tag}
				</Badge>
			</button>
		{/each}
		{#if currentTag}
			<button
				class=""
				onclick={(e: Event) => {
					e.stopPropagation();
					reset();
				}}
			>
				clear
			</button>
		{/if}
	</div>
</div>
