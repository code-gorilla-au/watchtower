<script lang="ts">
	import { OrgService } from "$lib/watchtower";
	import { Button } from "$components/ui/button";
	import { BaseInput } from "$components/base_input/index.js";

	const orgs = new OrgService();
	let state = $state({
		name: "",
		owner: ""
	});
	async function onSubmit(e: Event) {
		e.preventDefault();

		await orgs.create(state.name, state.owner);
	}
</script>

<h1 class="text-4xl">Register</h1>
<form onsubmit={onSubmit} class="flex flex-col">
	<div class="mb-4">Create new org</div>
	<BaseInput
		class=""
		id="friendly-name"
		label="Organisation name"
		description="Organisation's friendly name"
		bind:value={state.name}
	/>
	<BaseInput id="namespace" label="Github owner" description="Github's unique identifier" />

	<Button type="submit">Submit</Button>
	<div>
		{JSON.stringify(state)}
	</div>
</form>
