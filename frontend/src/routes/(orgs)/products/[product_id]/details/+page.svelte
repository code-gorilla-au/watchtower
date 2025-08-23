<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { RefreshCw, Pencil } from "@lucide/svelte";
	import { type PageProps } from "./$types";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid";
	import { RepoCard } from "$components/products";
	import { Button } from "$components/ui/button";
	import { PRGrid, SecurityGrid } from "$components/products/index.js";

	let { data }: PageProps = $props();

	let product = $derived(data.product);
	let repos = $derived(data.repos);
	let prs = $derived(data.prs);
	let securities = $derived(data.securities);

	async function syncProduct(e: Event) {
		e.preventDefault();

		await goto(`/products/${product.id}/sync`);
	}

	async function editProduct(e: Event) {
		e.preventDefault();
		await goto(`/products/${product.id}/edit`);
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
		<Button onclick={editProduct} variant="ghost" size="icon">
			<Pencil />
		</Button>
	</PageTitle>

	<div>
		<h2 class="text-xl text-muted-foreground">Security Vulnerabilities</h2>
		<SecurityGrid {securities} />
	</div>
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
