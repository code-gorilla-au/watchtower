<script lang="ts">
	import { watchtower } from "$lib/wailsjs/go/models";
	import { formatDate } from "$design/formats";
	import { Card, CardContent, CardTitle } from "$components/ui/card";
	import { Badge } from "$components/ui/badge";
	import { RefreshCw, Trash } from "@lucide/svelte";
	import { Button } from "$components/ui/button";

	import { goto } from "$app/navigation";
	import { productSvc } from "$lib/watchtower";

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
		<CardTitle class="flex items-center justify-between px-2">
			<span>{product.name}</span>
			<div>
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
			</div>
		</CardTitle>
		<CardContent>
			<div class="mb-2 flex justify-between text-sm">
				<p class="text-muted-foreground">Last updated:</p>
				<p>{formatDate(product.updated_at)}</p>
			</div>
			{#each product?.tags ?? [] as tag (tag)}
				<Badge variant="secondary" class="">{tag}</Badge>
			{/each}
		</CardContent>
	</Card>
</a>
