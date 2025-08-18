<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { type PageProps } from "./$types";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid";
	import { RepoCard } from "$components/repo_card";
	import { PRCard } from "$components/products";
	import { EmptySlate } from "$components/empty_slate/index.js";

	let { data }: PageProps = $props();

	let product = $derived(data.product);
	let repos = $derived(data.repos);
	let prs = $derived(data.prs);
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(`/products`);
		}}
		title="Product Details"
		subtitle="Details for product {product.name}"
	></PageTitle>

	<div>
		<h2 class="text-xl text-muted-foreground">Pull Requests</h2>
	</div>
	{#if prs.length === 0}
		<EmptySlate title="No Open PRs"></EmptySlate>
	{:else}
		<Grid>
			{#each prs as pr (pr.id)}
				<PRCard {pr} />
			{/each}
		</Grid>
	{/if}

	<div class="my-4">
		<h2 class="text-xl text-muted-foreground">Repositories</h2>
	</div>
	<Grid>
		{#each repos as repo (repo.id)}
			<RepoCard {repo} />
		{/each}
	</Grid>
</div>
