<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { RefreshCw } from "@lucide/svelte";
	import { type PageProps } from "./$types";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid";
	import { RepoCard } from "$components/repo_card";
	import { PRCard } from "$components/products";
	import { EmptySlate } from "$components/empty_slate/index.js";
	import { Button } from "$components/ui/button";

	let { data }: PageProps = $props();

	let product = $derived(data.product);
	let repos = $derived(data.repos);
	let prs = $derived(data.prs);

	async function syncProduct(e: Event) {
		e.preventDefault();

		await goto(`/products/${product.id}/sync`);
	}
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(`/products`);
		}}
		title="Product Details"
		subtitle="Details for product {product.name}"
	>
		<Button onclick={syncProduct} variant="ghost" size="icon">
			<RefreshCw />
		</Button>
	</PageTitle>

	<div class="my-4">
		<h2 class="text-xl text-muted-foreground">Pull Requests</h2>
	</div>
	{#if prs.length === 0}
		<EmptySlate title="No open PRs"></EmptySlate>
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
