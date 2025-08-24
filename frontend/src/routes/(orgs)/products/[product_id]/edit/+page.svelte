<script lang="ts">
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { type PageProps } from "./$types";

	import { ProductUpdateForm, type ProductUpdateFormData } from "$components/products";
	import { productSvc } from "$lib/watchtower";

	let { data }: PageProps = $props();

	let product = $derived(data.product);

	type PageState = {
		error: string | undefined;
		loading: boolean;
	};

	let pageState = $state<PageState>({
		error: undefined,
		loading: false
	});

	async function onUpdate(formData: ProductUpdateFormData) {
		try {
			pageState.loading = true;
			pageState.error = undefined;
			await productSvc.update(product.id, formData.name, formData.tags.split(","));
		} catch (e) {
			const err = e as Error;
			pageState.error = err.message;
		} finally {
			pageState.loading = false;
		}
	}

	async function goBack(e: Event) {
		e.preventDefault();
		await goto("/products");
	}
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto("/products");
		}}
		title="Edit Product"
		subtitle="Edit {product.name} details"
	/>

	<ProductUpdateForm
		mode="update"
		error={pageState.error}
		loading={pageState.loading}
		{product}
		{onUpdate}
		onCancel={goBack}
	/>
</div>
