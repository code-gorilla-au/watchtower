<script lang="ts">
	import * as Accordion from "$components/ui/accordion";
	import { GridHeader } from "$components/grid";
	import { SecurityGrid } from "$components/products";
	import { products } from "$lib/wailsjs/go/models";
	import { TagsFilter } from "$lib/hooks/filters.svelte";

	type Props = {
		securities: products.SecurityDTO[];
	};

	let { securities }: Props = $props();

	const tagsFilter = new TagsFilter(securities, "tag");
</script>

<Accordion.Item value="security">
	<Accordion.Trigger class="text-left">
		<GridHeader
			dataLength={securities.length}
			tags={tagsFilter.tags}
			currentTag={tagsFilter.currentTag}
			filterByTag={(tag) => {
				tagsFilter.filterByTag(tag);
			}}
			reset={() => {
				tagsFilter.reset();
			}}
			title="Security Vulnerabilities"
		/>
	</Accordion.Trigger>
	<Accordion.Content class="mb-5">
		<SecurityGrid securities={tagsFilter.data} />
	</Accordion.Content>
</Accordion.Item>
