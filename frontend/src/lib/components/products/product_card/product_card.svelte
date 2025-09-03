<script lang="ts">
	import { watchtower } from "$lib/wailsjs/go/models";
	import { formatDate } from "$lib/hooks/formats";
	import { Card, CardContent, CardTitle } from "$components/ui/card";
	import { Badge } from "$components/ui/badge";
	import { RefreshCw, Trash } from "@lucide/svelte";
	import { Button } from "$components/ui/button";
	import { goto } from "$app/navigation";
	import { CardAction, CardHeader } from "$components/ui/card/index.js";

	type Props = {
		product: watchtower.ProductDTO;
		onDelete?: () => void;
	};

	let { product, onDelete }: Props = $props();

	async function syncProduct(id: number) {
		await goto(`/products/${id}/sync`);
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
						onDelete?.();
					}}
					size="icon"
					variant="ghost"
				>
					<Trash />
				</Button>
			</CardAction>
		</CardHeader>

		<CardContent>
			<p>{product.description}</p>
			<div class="my-2 flex justify-between text-sm">
				<p class="text-muted-foreground">Last updated:</p>
				<p>{formatDate(product.updated_at)}</p>
			</div>
			{#each product?.tags ?? [] as tag (tag)}
				<Badge class="">{tag}</Badge>
			{/each}
		</CardContent>
	</Card>
</a>
