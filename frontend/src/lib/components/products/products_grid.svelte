<script lang="ts">
	import { resolve } from "$app/paths";
	import { EmptySlate } from "$components/empty_slate";
	import { ProductCard } from "$components/products/product_card";
	import { Grid } from "$components/grid";
	import { products as p } from "$lib/wailsjs/go/models";
	import { productSvc } from "$lib/watchtower";
	import { invalidateAll } from "$app/navigation";
	import { cn } from "$lib/utils";
	import { toast } from "svelte-sonner";

	type Props = {
		products: p.ProductDTO[];
		class?: string;
	};

	let { products, class: className }: Props = $props();

	function findProductById(id: number) {
		return products.find((p) => p.id === id);
	}

	async function deleteProduct(id: number) {
		try {
			const product = findProductById(id);

			await productSvc.delete(id);

			toast.success(`Product ${product?.name} deleted`, {
				position: "top-right"
			});

			await invalidateAll();
		} catch (e) {
			console.error(e);
		}
	}
</script>

<div class={cn(className)}>
	{#if products.length === 0}
		<EmptySlate caution={true} title="No products">
			<a href={resolve("/products/create")} class="text-xs text-muted-foreground underline">
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
