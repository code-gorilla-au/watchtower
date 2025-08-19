<script lang="ts">
	import type { PageProps } from "./$types";
	import { productSvc } from "$lib/watchtower";
	import { ProductUpdateForm, type ProductUpdateFormData } from "$components/products/index.js";
	import { goto } from "$app/navigation";

	let { data }: PageProps = $props();

	let organisation = $derived(data.organisation);

	type PageState = {
		error: string | undefined;
		loading: boolean;
	};

	let pageState = $state<PageState>({
		error: undefined,
		loading: false
	});

	async function onSubmit(formData: ProductUpdateFormData) {
		try {
			await productSvc.create(
				formData.name,
				formData.description,
				organisation.id,
				formData.tags.split(",")
			);

			await goto("/");
		} catch (e) {
			const err = e as Error;
			pageState.error = err.message;
		} finally {
			pageState.loading = false;
		}
	}
</script>

<div class="p-3">
	<h1 class="text-4xl">Register - Product</h1>
	<div class="mb-10">Add a product to an organisation {organisation.friendly_name}</div>
	<ProductUpdateForm
		mode="create"
		error={pageState.error}
		loading={pageState.loading}
		onCreate={onSubmit}
	/>
</div>
