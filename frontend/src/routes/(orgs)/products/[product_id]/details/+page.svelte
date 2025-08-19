<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { RefreshCw } from "@lucide/svelte";
	import { type PageProps } from "./$types";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid";
	import { RepoCard } from "$components/products";
	import { Button } from "$components/ui/button";
	import { PRGrid } from "$components/products/index.js";

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
	<PRGrid {prs} />

	<div class="my-4">
		<h2 class="text-xl text-muted-foreground">Repositories</h2>
	</div>
	<Grid>
		{#each repos as repo (repo.id)}
			<RepoCard {repo} />
		{/each}
	</Grid>
</div>
