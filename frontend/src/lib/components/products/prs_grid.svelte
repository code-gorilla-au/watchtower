<script lang="ts">
	import { PRCard } from "$components/products/index";
	import { EmptySlate } from "$components/empty_slate";
	import { Grid } from "$components/grid";
	import { products } from "$lib/wailsjs/go/models";

	type Props = {
		prs: products.PullRequestDTO[];
	};
	let { prs }: Props = $props();
</script>

{#if prs.length === 0}
	<EmptySlate title="No open PRs"></EmptySlate>
{:else}
	<Grid>
		{#each prs as pr (pr.id)}
			<PRCard {pr} />
		{/each}
	</Grid>
{/if}
