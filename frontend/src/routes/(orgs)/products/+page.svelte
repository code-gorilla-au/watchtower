<script lang="ts">
	import { type PageProps } from "./$types";

	let { data }: PageProps = $props();
	const products = $derived(data.products);
</script>

<h1 class="text-4xl">Products</h1>

{#each products as product (product.id)}
	<div>{product.name}</div>
{/each}
