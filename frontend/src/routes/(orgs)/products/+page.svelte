<script lang="ts">
	import { type PageProps } from "./$types";
	import { Plus } from "@lucide/svelte";
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { Button } from "$components/ui/button";
	import { ProductsGrid } from "$components/products/index.js";
	import { resolve } from "$app/paths";

	let { data }: PageProps = $props();
	const products = $derived(data.products ?? []);
	const organisation = $derived(data.organisation);

	async function createProduct() {
		await goto(resolve("/products/create"));
	}
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(resolve("/"));
		}}
		title="Products"
		subtitle="All products for {organisation?.friendly_name}"
	>
		<Button onclick={createProduct}><Plus /></Button>
	</PageTitle>
	<ProductsGrid {products} />
</div>
