<script lang="ts">
	import { PageTitle } from "$components/page_title/index.js";
	import { BaseInput } from "$components/base_input/index.js";
	import { Button } from "$components/ui/button";
	import type { PageProps } from "./$types";
	import { Switch } from "$components/ui/switch";
	import { Label } from "$components/ui/label";
	import { orgSvc } from "$lib/watchtower";
	import { goto } from "$app/navigation";

	let { data }: PageProps = $props();

	type FormData = {
		id: number;
		friendly_name: string;
		namespace: string;
		default_org: boolean;
	};

	const formData = $state<FormData>({
		id: data.org?.id,
		friendly_name: data.org.friendly_name,
		namespace: data.org.namespace,
		default_org: data.org.default_org
	});

	async function updateOrg(e: Event) {
		e.preventDefault();

		await orgSvc.update({
			id: formData.id,
			friendlyName: formData.friendly_name,
			owner: formData.namespace,
			defaultOrg: formData.default_org
		});

		await goto("/dashboard");
	}
</script>

<div class="w-full p-2">
	<PageTitle title="Organisation"></PageTitle>

	<form method="POST" onsubmit={updateOrg}>
		<BaseInput
			id="friendly-name"
			label="Name"
			description="Organisation's friendly name"
			bind:value={formData.friendly_name}
		/>
		<BaseInput
			id="namespace"
			label="Github owner"
			description="Github's unique identifier"
			bind:value={formData.namespace}
		/>
		<div class="flex justify-between">
			<Label for="default-org">Default organisation</Label>
			<Switch id="default-org" bind:checked={formData.default_org} />
		</div>
		<div class="my-10 flex justify-end">
			<Button type="submit">Update</Button>
		</div>
	</form>
</div>
