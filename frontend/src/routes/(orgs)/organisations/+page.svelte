<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { goto } from "$app/navigation";
	import type { PageProps } from "./$types";
	import { Button } from "$components/ui/button";
	import { Plus } from "@lucide/svelte";
	import { resolve } from "$app/paths";
	import { OrgsGrid } from "$components/orgs/index.js";

	let { data }: PageProps = $props();

	const orgs = $derived(data.orgs);

	async function createNewOrg(e: Event) {
		e.preventDefault();

		await goto(resolve("/organisations/create"));
	}
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(resolve("/"));
		}}
		title="Organisations"
		subtitle="Manage your organisations"
	>
		<Button onclick={createNewOrg}><Plus /></Button>
	</PageTitle>

	<OrgsGrid {orgs} />
</div>
