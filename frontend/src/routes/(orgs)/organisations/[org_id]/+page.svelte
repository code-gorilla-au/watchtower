<script lang="ts">
	import { PageTitle } from "$components/page_title/index.js";
	import type { PageProps } from "./$types";
	import { orgSvc } from "$lib/watchtower";
	import { goto } from "$app/navigation";
	import { OrgUpdateForm, type OrgUpdateFormData } from "$components/orgs";
	import { resolve } from "$app/paths";

	let { data }: PageProps = $props();
	let org = $derived(data.org);

	async function updateOrg(formData: OrgUpdateFormData) {
		await orgSvc.update({
			id: formData.id,
			friendlyName: formData.friendly_name,
			owner: formData.namespace,
			defaultOrg: formData.default_org,
			description: formData.description
		});

		await goto(resolve("/"));
	}

	async function setDefault() {
		await orgSvc.setDefault(org.id);
		await goto(resolve("/"));
	}
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(resolve("/"));
		}}
		title="Organisation"
	/>

	<OrgUpdateForm {org} mode="update" onUpdate={updateOrg} onSetDefault={setDefault} />
</div>
