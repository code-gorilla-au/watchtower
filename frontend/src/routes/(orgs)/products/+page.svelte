<script lang="ts">
	import { type PageProps } from "./$types";
	import { Plus } from "@lucide/svelte";
	import { Grid } from "$components/grid";
	import { EmptySlate } from "$components/empty_slate";
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { Button } from "$components/ui/button";
	import { Card } from "$components/ui/card";
	import { CardContent, CardTitle } from "$components/ui/card/index.js";
	import { Badge } from "$components/ui/badge/index.js";
	import { RefreshCw } from "@lucide/svelte";
	import { formatDate } from "$design/formats";

	let { data }: PageProps = $props();
	const products = $derived(data.products ?? []);
	const organisation = $derived(data.organisation);

	async function createProduct() {
		await goto("/products/create");
	}
	async function syncProduct(id: number) {
		await goto(`/products/sync/${id}`);
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
					<Card class="w-full">
						<CardTitle class="flex items-center justify-between px-2">
							<span>{product.name}</span>
							<Button
								onclick={async (e: Event) => {
									e.preventDefault();
									await syncProduct(product.id);
								}}
								size="icon"
								variant="ghost"><RefreshCw /></Button
							>
						</CardTitle>
						<CardContent>
							<div class="mb-2 flex justify-between text-sm">
								<p class="text-muted-foreground">Last updated:</p>
								<p>{formatDate(product.updated_at)}</p>
							</div>
							{#each product?.tags?.split(",") ?? [] as tag (tag)}
								<Badge variant="secondary" class="">{tag}</Badge>
							{/each}
						</CardContent>
					</Card>
				{/each}
			</Grid>
		{/if}
	</div>
</div>
