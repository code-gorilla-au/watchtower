<script lang="ts">
	import { type PageProps } from "./$types";
	import { Grid } from "$components/grid";
	import { EmptySlate } from "$components/empty_slate";

	let { data }: PageProps = $props();
	const products = $derived(data.products ?? []);
</script>

<div class="w-full p-2">
	<h1 class="text-4xl">Products</h1>
	<div class="p-2">
		{#if products.length === 0}
			<EmptySlate caution={true} title="No products">
				<a href="/products/create" class="text-xs text-muted-foreground underline"
					>Create a product to get started</a
				>
			</EmptySlate>
		{:else}
			<Grid>
				{#each products as product (product.id)}
					<div>{product.name}</div>
				{/each}
			</Grid>
		{/if}
	</div>
</div>
