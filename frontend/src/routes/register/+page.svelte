<script lang="ts">
	import { OrgService } from "$lib/watchtower";
	import { Button } from "$components/ui/button";

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
	<label for="org-name">
		Org friendly name
		<input required id="org-name" bind:value={state.name} />
	</label>
	<label for="org-name">
		Org name space (Github owner)
		<input required id="org-owner" bind:value={state.owner} />
	</label>

	<Button type="submit">Submit</Button>
	<div>
		{JSON.stringify(state)}
	</div>
</form>
