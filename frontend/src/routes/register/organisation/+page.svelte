<script lang="ts">
	import { OrgService } from "$lib/watchtower";
	import { Button } from "$components/ui/button";
	import { BaseInput } from "$components/base_input/index.js";
	import { goto } from "$app/navigation";

	const orgs = new OrgService();

	type FormData = {
		name: string;
		owner: string;
		patToken: string;
		error: string | null;
	};

	let state = $state<FormData>({
		name: "",
		owner: "",
		patToken: "",
		error: null
	});
	async function onSubmit(e: Event) {
		try {
			e.preventDefault();
			state.error = null;

			await orgs.create(state.name, state.owner);
			await goto("/register/product");
		} catch (e) {
			const err = e as Error;
			state.error = err.message;
		}
	}
</script>

<div class="p-3">
	<h1 class="text-4xl">Register - Organisation</h1>
	<div class="mb-10">Create a new organisation</div>
	<form
		method="POST"
		onsubmit={onSubmit}
		class="mx-auto flex max-w-sm flex-col items-center justify-center"
	>
		<BaseInput
			class=""
			id="friendly-name"
			label="Organisation name"
			description="Organisation's friendly name"
			bind:value={state.name}
		/>

		<BaseInput
			id="namespace"
			label="Github owner"
			description="Github's unique identifier"
			bind:value={state.owner}
		/>

		<BaseInput
			id="pat-token"
			label="Personal access token"
			description=" Readonly personal access token for the organisation"
			bind:value={state.patToken}
		/>

		<div class="my-10 flex w-full justify-end">
			<Button type="submit">Add product</Button>
		</div>
	</form>
</div>
