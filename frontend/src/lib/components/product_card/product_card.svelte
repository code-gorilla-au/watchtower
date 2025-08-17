<script lang="ts">
	import { watchtower } from "$lib/wailsjs/go/models";
	import { formatDate } from "$design/formats";
	import { Card, CardContent, CardTitle } from "$components/ui/card";
	import { Badge } from "$components/ui/badge";
	import { RefreshCw, Trash } from "@lucide/svelte";
	import { Button } from "$components/ui/button";

	import { goto } from "$app/navigation";
	import { productSvc } from "$lib/watchtower";
	import { CardAction, CardHeader } from "$components/ui/card/index.js";

	type Props = {
		product: watchtower.ProductDTO;
	};

	let { product }: Props = $props();

	async function syncProduct(id: number) {
		await goto(`/products/${id}/sync`);
	}

	async function deleteProduct(id: number) {
		await productSvc.delete(id);
		await goto("/dashboard");
	}
</script>

<a href={`/products/${product.id}/details`}>
	<Card class="w-full cursor-pointer hover:bg-muted/30">
		<CardHeader class="flex items-center justify-between">
			<CardTitle>
				{product.name}
			</CardTitle>
			<CardAction>
				<Button
					onclick={async (e: Event) => {
						e.preventDefault();
						await syncProduct(product.id);
					}}
					size="icon"
					variant="ghost"
				>
					<RefreshCw />
				</Button>
				<Button
					onclick={async (e: Event) => {
						e.preventDefault();
						await deleteProduct(product.id);
					}}
					size="icon"
					variant="ghost"
				>
					<Trash />
				</Button>
			</CardAction>
		</CardHeader>

		<CardContent>
			<div class="mb-2 flex justify-between text-sm">
				<p class="text-muted-foreground">Last updated:</p>
				<p>{formatDate(product.updated_at)}</p>
			</div>
			{#each product?.tags ?? [] as tag (tag)}
				<Badge class="">{tag}</Badge>
			{/each}
		</CardContent>
	</Card>
</a>
