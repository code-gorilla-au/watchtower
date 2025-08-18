<script lang="ts">
	import { type PageProps } from "./$types";
	import { Plus } from "@lucide/svelte";
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { Button } from "$components/ui/button";
	import { ProductsGrid } from "$components/products/index.js";

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
	<ProductsGrid {products} />
</div>
