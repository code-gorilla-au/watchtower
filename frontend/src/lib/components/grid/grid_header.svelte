<script lang="ts" generics="T extends object">
	import { Badge } from "$components/ui/badge";
	import { type FilterTag, type FilterTagValue, TagsFilter } from "$lib/hooks/filters.svelte";

	type Props<T> = {
		title: string;
		data: T[];
		tagField: FilterTag<T>;
	};

	let { data, tagField, title }: Props<T> = $props();

	const tagsFilter = new TagsFilter(data, tagField);

	function styleCurrentTag(tag: FilterTagValue<T>, currentTag?: FilterTagValue<T>) {
		if (currentTag === tag) {
			return "default";
		}

		return "outline";
	}
</script>

<div class="">
	<h2 class="mb-1 text-xl text-muted-foreground">{title} ({data.length})</h2>
	<div class="flex gap-2">
		{#each tagsFilter.tags as tag (tag)}
			<button
				onclick={(e: Event) => {
					e.stopPropagation();
					tagsFilter.filterByTag(tag);
				}}
			>
				<Badge variant={styleCurrentTag(tag, tagsFilter?.currentTag)} class="text-xs">
					{tag}
				</Badge>
			</button>
		{/each}
		{#if tagsFilter.currentTag}
			<button
				class=""
				onclick={(e: Event) => {
					e.stopPropagation();
					tagsFilter.reset();
				}}
			>
				clear
			</button>
		{/if}
	</div>
</div>
