<script lang="ts">
	import type { PageProps } from "./$types";
	import { productSvc } from "$lib/watchtower";
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { ProductUpdateForm, type ProductUpdateFormData } from "$components/products/index.js";

	let { data }: PageProps = $props();

	let organisation = $derived(data.org);

	type FormData = {
		error?: string;
		loading: boolean;
	};

	const form = $state<FormData>({
		error: undefined,
		loading: false
	});

	async function onSubmit(formData: ProductUpdateFormData) {
		try {
			if (!organisation) {
				form.error = "Organisation not found";
				return;
			}

			const product = await productSvc.create(
				formData.name,
				formData.description,
				organisation?.id,
				formData.tags.split(",")
			);
			await goto(`/products/${product.id}/sync`);
			return;
		} catch (e) {
			const err = e as Error;
			form.error = err.message;
		}
	}
</script>

<div class="w-full p-2">
	<PageTitle
		title="Add product"
		backAction={async () => {
			await goto("/products");
		}}
		subtitle="Add a product to an organisation {organisation?.friendly_name}"
	/>

	<ProductUpdateForm mode="create" error={form.error} loading={form.loading} onCreate={onSubmit} />
</div>
