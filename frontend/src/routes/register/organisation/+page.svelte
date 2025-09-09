<script lang="ts">
	import { goto } from "$app/navigation";
	import { OrgUpdateForm, type OrgUpdateFormData } from "$components/orgs";
	import { orgSvc } from "$lib/watchtower";
	import { PageTitle } from "$components/page_title/index.js";
	import { resolve } from "$app/paths";

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

			await goto(resolve("/register/product"));
		} catch (e) {
			const err = e as Error;
			pageState.error = err.message;
		} finally {
			pageState.loading = false;
		}
	}
</script>

<div class="page-container">
	<PageTitle class="mb-10" title="Register - Organisation" subtitle="Create a new organisation" />
	<OrgUpdateForm
		loading={pageState.loading}
		error={pageState.error}
		mode="create"
		onCreate={onSubmit}
	/>
</div>
