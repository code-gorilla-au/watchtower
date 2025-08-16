<script lang="ts">
	import { PageTitle } from "$components/page_title/index.js";
	import { OrgUpdateForm, type OrgUpdateFormData } from "$components/orgs/index.js";
	import { orgSvc } from "$lib/watchtower";
	import { goto } from "$app/navigation";

	async function createOrg(formData: OrgUpdateFormData) {
		await orgSvc.create(formData.friendly_name, formData.namespace, formData.token);
		await goto("/");
	}
</script>

<div class="w-full p-2">
	<PageTitle title="Create Organisation" subtitle="Create a new organisation" />

	<OrgUpdateForm mode="create" onCreate={createOrg} />
</div>
