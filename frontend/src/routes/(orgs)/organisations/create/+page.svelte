<script lang="ts">
	import { toast } from "svelte-sonner";
	import { PageTitle } from "$components/page_title/index.js";
	import { OrgUpdateForm, type OrgUpdateFormData } from "$components/orgs/index.js";
	import { orgSvc } from "$lib/watchtower";
	import { goto } from "$app/navigation";
	import { resolve } from "$app/paths";

	async function createOrg(formData: OrgUpdateFormData) {
		await orgSvc.create(
			formData.friendly_name,
			formData.namespace,
			formData.token,
			formData.description
		);

		toast.success("Organisation created", {
			position: "top-right"
		});

		await goto(resolve("/"));
	}
</script>

<div class="page-container">
	<PageTitle title="Create Organisation" subtitle="Create a new organisation" />

	<OrgUpdateForm
		mode="create"
		onCreate={createOrg}
		onCancel={() => goto(resolve("/organisations"))}
	/>
</div>
