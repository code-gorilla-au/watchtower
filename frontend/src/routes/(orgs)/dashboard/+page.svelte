<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";
	import { PRGrid, ProductsGrid } from "$components/products";

	let { data }: PageProps = $props();
	let org = $derived(data.organisation);
	let products = $derived(data.products);
	let prs = $derived(data.prs);
</script>

<div class="page-container">
	<PageTitle title="Dashboard" subtitle={org?.friendly_name} />

	<div class="my-4">
		<h3 class="text-xl text-muted-foreground">Pull Requests</h3>
	</div>
	<PRGrid {prs} />

	<div class="my-4">
		<h3 class="text-xl text-muted-foreground">Products</h3>
	</div>
	<ProductsGrid {products} />
</div>
