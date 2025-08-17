<script lang="ts">
	import { goto } from "$app/navigation";
	import { OrgUpdateForm, type OrgUpdateFormData } from "$components/orgs";
	import { orgSvc } from "$lib/watchtower";

	type PageState = {
		error: string | undefined;
		loading: boolean;
	};

	const pageState = $state<PageState>({
		error: undefined,
		loading: false
	});

	async function onSubmit(formData: OrgUpdateFormData) {
		try {
			pageState.loading = true;
			pageState.error = undefined;

			await orgSvc.create(
				formData.friendly_name,
				formData.namespace,
				formData.token,
				formData.description
			);

			await goto("/register/product");
		} catch (e) {
			const err = e as Error;
			pageState.error = err.message;
		} finally {
			pageState.loading = false;
		}
	}
</script>

<div class="p-3">
	<h1 class="text-4xl">Register - Organisation</h1>
	<div class="mb-10">Create a new organisation</div>
	<OrgUpdateForm
		loading={pageState.loading}
		error={pageState.error}
		mode="create"
		onCreate={onSubmit}
	/>
</div>
