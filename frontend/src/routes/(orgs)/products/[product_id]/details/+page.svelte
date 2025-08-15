<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { type PageProps } from "./$types";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid";

	let { data }: PageProps = $props();

	let product = $derived(data.product);
	let repos = $derived(data.repos);
</script>

<div class="w-full p-2">
	<PageTitle
		backAction={async () => {
			await goto(`/products`);
		}}
		title="Product Details"
		subtitle="Details for product {product.name}"
	></PageTitle>

	<Grid>
		{#each repos as repo (repo.id)}
			{JSON.stringify(repo)}
		{/each}
	</Grid>
</div>
