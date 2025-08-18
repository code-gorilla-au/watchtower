<script lang="ts">
	import { EmptySlate } from "$components/empty_slate";
	import { ProductCard } from "$components/products/product_card";
	import { Grid } from "$components/grid";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { productSvc } from "$lib/watchtower";
	import { invalidateAll } from "$app/navigation";
	import { cn } from "$lib/utils";

	type Props = {
		products: watchtower.ProductDTO[];
		class?: string;
	};

	let { products, class: className }: Props = $props();

	async function deleteProduct(id: number) {
		try {
			await productSvc.delete(id);
			await invalidateAll();
		} catch (e) {
			console.error(e);
		}
	}
</script>

<div class={cn("p-2", className)}>
	{#if products.length === 0}
		<EmptySlate caution={true} title="No products">
			<a href="/products/create" class="text-xs text-muted-foreground underline">
				Create a product to get started
			</a>
		</EmptySlate>
	{:else}
		<Grid>
			{#each products as product (product.id)}
				<ProductCard {product} onDelete={() => deleteProduct(product.id)} />
			{/each}
		</Grid>
	{/if}
</div>
