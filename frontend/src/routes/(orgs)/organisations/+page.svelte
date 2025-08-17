<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid/index.js";
	import type { PageProps } from "./$types";
	import { OrgCard } from "$components/orgs";
	import { Button } from "$components/ui/button";
	import { Plus } from "@lucide/svelte";

	let { data }: PageProps = $props();

	const orgs = $derived(data.orgs);

	async function createNewOrg(e: Event) {
		e.preventDefault();

		await goto("/organisations/create");
	}
</script>

<div class="w-full p-2">
	<PageTitle
		backAction={async () => {
			await goto("/");
		}}
		title="Organisations"
		subtitle="Manage your organisations"
	>
		<Button onclick={createNewOrg}><Plus /></Button>
	</PageTitle>

	<Grid>
		{#each orgs as org (org.id)}
			<OrgCard {org} />
		{/each}
	</Grid>
</div>
