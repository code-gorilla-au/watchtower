<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { type PageProps } from "./$types";
	import { goto } from "$app/navigation";

	let { data }: PageProps = $props();

	let product = $derived(data.product);
</script>

<div class="w-full p-2">
	<PageTitle
		backAction={async () => {
			await goto(`/products`);
		}}
		title="Product Details"
		subtitle="Details for product {product.name}"
	></PageTitle>
</div>
