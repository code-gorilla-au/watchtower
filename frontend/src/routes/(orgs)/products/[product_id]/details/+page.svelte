<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { type PageProps } from "./$types";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid";
	import { RepoCard } from "$components/repo_card";

	let { data }: PageProps = $props();

	let product = $derived(data.product);
	let repos = $derived(data.repos);
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(`/products`);
		}}
		title="Product Details"
		subtitle="Details for product {product.name}"
	></PageTitle>

	<Grid>
		{#each repos as repo (repo.id)}
			<RepoCard {repo} />
		{/each}
	</Grid>
</div>
