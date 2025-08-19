<script lang="ts">
	import { BaseInput } from "$components/base_input";
	import { Button } from "$components/ui/button";
	import { LoaderCircle } from "@lucide/svelte";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { type ProductUpdateFormData } from "./types";

	type Props = {
		product?: watchtower.ProductDTO;
		mode: "create" | "update";
		error?: string;
		loading?: boolean;
		onUpdate?: (formData: ProductUpdateFormData) => void;
		onCreate?: (formData: ProductUpdateFormData) => void;
	};

	let { product, mode, onUpdate, onCreate, error, loading }: Props = $props();

	const formData = $state<ProductUpdateFormData>({
		id: product?.id ?? 0,
		name: product?.name ?? "",
		description: product?.description ?? "",
		tags: product?.tags.join(",") ?? ""
	});

	function handleFormAction(e: Event) {
		e.preventDefault();

		if (mode === "create") {
			onCreate?.(formData);
		} else {
			onUpdate?.(formData);
		}
	}
</script>

<form
	method="POST"
	onsubmit={handleFormAction}
	class="mx-auto flex max-w-sm flex-col items-center justify-center"
>
	<BaseInput
		id="name"
		label="Name"
		required
		description="Product's friendly name"
		bind:value={formData.name}
	/>

	<BaseInput
		id="description"
		label="Description"
		required
		description="Short description of the product"
		bind:value={formData.description}
	/>

	<BaseInput
		id="tags"
		label="Tags"
		required
		description="Comma separated list of tags for the product"
		bind:value={formData.tags}
	/>

	<div class="my-10 flex w-full justify-end gap-3">
		{#if error}
			<div class="border-destructive text-destructive">
				<p>{error}</p>
			</div>
		{/if}

		<Button type="submit" class="capitalize">
			{#if loading}
				<LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
			{/if}
			{mode}
		</Button>
	</div>
</form>
