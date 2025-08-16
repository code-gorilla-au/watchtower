<script lang="ts">
	import { BaseInput } from "$components/base_input";
	import { Label } from "$components/ui/label";
	import { Switch } from "$components/ui/switch";
	import { Button } from "$components/ui/button";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { type OrgUpdateFormData } from "./types";

	type Props = {
		org?: watchtower.OrganisationDTO;
		mode: "create" | "update";
		onUpdate?: (formData: OrgUpdateFormData) => void;
		onCreate?: (formData: OrgUpdateFormData) => void;
		onSetDefault?: () => void;
	};

	let { org, mode, onUpdate, onCreate, onSetDefault }: Props = $props();

	const formData = $state<OrgUpdateFormData>({
		id: org?.id ?? 0,
		friendly_name: org?.friendly_name ?? "",
		namespace: org?.namespace ?? "",
		token: "",
		default_org: org?.default_org ?? true
	});

	function handleFormAction(e: Event) {
		e.preventDefault();

		if (mode === "create") {
			onCreate?.(formData);
		} else {
			onUpdate?.(formData);
		}
	}
</script>

<form
	method="POST"
	onsubmit={handleFormAction}
	class="mx-auto flex max-w-sm flex-col items-center justify-center"
>
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
	{#if mode === "create"}
		<BaseInput
			id="token"
			label="Github token"
			description="Github's personal access token"
			bind:value={formData.token}
		/>
	{/if}
	<div class="my-3 flex w-full justify-between">
		<Label for="default-org">Default organisation</Label>
		<Switch id="default-org" bind:checked={formData.default_org} />
	</div>
	<div class="my-10 flex w-full justify-end gap-3">
		{#if mode === "update"}
			<Button onclick={onSetDefault} variant="outline">Set default</Button>
		{/if}

		<Button type="submit" class="capitalize">{mode}</Button>
	</div>
</form>
