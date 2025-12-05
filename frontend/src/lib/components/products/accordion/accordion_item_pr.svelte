<script lang="ts">
	import * as Accordion from "$components/ui/accordion";
	import { GridHeader } from "$components/grid";
	import { PRGrid } from "$components/products";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { type FilterTagValue, TagsFilter } from "$lib/hooks/filters.svelte";

	type Props = {
		prs: watchtower.PullRequestDTO[];
		fieldTag?: FilterTagValue<watchtower.PullRequestDTO>;
	};

	let { prs, fieldTag = "tag" }: Props = $props();

	const tagsFilter = new TagsFilter(prs, fieldTag);
</script>

<Accordion.Item value="prs">
	<Accordion.Trigger class="text-left">
		<GridHeader
			dataLength={prs.length}
			tags={tagsFilter.tags}
			currentTag={tagsFilter.currentTag}
			filterByTag={(tag) => tagsFilter.filterByTag(tag)}
			reset={() => tagsFilter.reset()}
			title="Pull Requests"
		/>
	</Accordion.Trigger>
	<Accordion.Content class="mb-5">
		<PRGrid prs={tagsFilter.data} />
	</Accordion.Content>
</Accordion.Item>
