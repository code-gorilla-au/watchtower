<script lang="ts">
	import { BaseInput } from "$components/base_input";
	import { Label } from "$components/ui/label";
	import { Switch } from "$components/ui/switch";
	import { Button } from "$components/ui/button";
	import { LoaderCircle } from "@lucide/svelte";
	import { organisations } from "$lib/wailsjs/go/models";
	import { type OrgUpdateFormData } from "./types";

	type Props = {
		org?: organisations.OrganisationDTO;
		mode: "create" | "update";
		error?: string;
		loading?: boolean;
		onUpdate?: (formData: OrgUpdateFormData) => void;
		onCreate?: (formData: OrgUpdateFormData) => void;
		onCancel?: () => void;
		onSetDefault?: () => void;
	};

	let { org, mode, onUpdate, onCreate, onSetDefault, error, loading, onCancel }: Props = $props();

	const formData = $state<OrgUpdateFormData>({
		id: org?.id ?? 0,
		friendly_name: org?.friendly_name ?? "",
		description: org?.description ?? "",
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
		required
		description="Organisation's friendly name"
		bind:value={formData.friendly_name}
	/>

	<BaseInput
		id="description"
		label="Description"
		required
		description="Short description of the organisation"
		bind:value={formData.description}
	/>

	<BaseInput
		id="namespace"
		label="Github owner"
		required
		description="Github's unique identifier"
		bind:value={formData.namespace}
	/>
	{#if mode === "create"}
		<BaseInput
			id="token"
			label="Github token"
			required
			description="Github's personal access token"
			bind:value={formData.token}
		/>
	{/if}
	<div class="my-3 flex w-full justify-between">
		<Label for="default-org">Default organisation</Label>
		<Switch required id="default-org" bind:checked={formData.default_org} />
	</div>

	{#if error}
		<div class="border-destructive text-destructive">
			<p>{error}</p>
		</div>
	{/if}
	<div class="my-10 flex w-full justify-end">
		{#if onCancel}
			<div class="w-full">
				<Button
					onclick={(e: Event) => {
						e.preventDefault();
						onCancel?.();
					}}
					variant="outline">Cancel</Button
				>
			</div>
		{/if}

		<div class="flex gap-4">
			{#if mode === "update"}
				<Button onclick={onSetDefault} variant="outline">Set default</Button>
			{/if}

			<Button type="submit" class="capitalize">
				{#if loading}
					<LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				{mode}
			</Button>
		</div>
	</div>
</form>
