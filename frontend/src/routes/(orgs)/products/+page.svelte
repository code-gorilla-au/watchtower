<script lang="ts">
	import { type PageProps } from "./$types";
	import { Plus } from "@lucide/svelte";
	import { Grid } from "$components/grid";
	import { EmptySlate } from "$components/empty_slate";
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { Button } from "$components/ui/button";
	import { productSvc } from "$lib/watchtower";
	import { ProductCard } from "$components/product_card";

	let { data }: PageProps = $props();
	const products = $derived(data.products ?? []);
	const organisation = $derived(data.organisation);

	async function createProduct() {
		await goto("/products/create");
	}
</script>

<div class="w-full p-2">
	<PageTitle
		backAction={async () => {
			await goto("/");
		}}
		title="Products"
		subtitle="All products for {organisation?.friendly_name}"
	>
		<Button onclick={createProduct}><Plus /></Button>
	</PageTitle>
	<div class="p-2">
		{#if products.length === 0}
			<EmptySlate caution={true} title="No products">
				<a href="/products/create" class="text-xs text-muted-foreground underline">
					Create a product to get started
				</a>
			</EmptySlate>
		{:else}
			<Grid>
				{#each products as product (product.id)}
					<ProductCard {product} />
				{/each}
			</Grid>
		{/if}
	</div>
</div>
