<script lang="ts">
	import type { PageProps } from "./$types";
	import { productSvc } from "$lib/watchtower";
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { ProductUpdateForm, type ProductUpdateFormData } from "$components/products/index.js";
	import { resolve } from "$app/paths";
	import { toast } from "svelte-sonner";

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

			toast.success("Product created", {
				position: "top-right"
			});

			await goto(resolve(`/products/${product.id}/sync`));
		} catch (e) {
			const err = e as Error;
			form.error = err.message;
		}
	}

	async function onCancel(e: Event) {
		e.preventDefault();
		await goto(resolve("/products"));
	}
</script>

<div class="page-container">
	<PageTitle
		title="Add product"
		backAction={async () => {
			await goto(resolve("/products"));
		}}
		subtitle="Add a product to an organisation {organisation?.friendly_name}"
	/>

	<ProductUpdateForm
		mode="create"
		error={form.error}
		loading={form.loading}
		onCreate={onSubmit}
		{onCancel}
	/>
</div>
