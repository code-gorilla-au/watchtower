<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid/index.js";
	import type { PageProps } from "./$types";
	import { OrgCard } from "$components/orgs";
	import { Button } from "$components/ui/button";
	import { Plus } from "@lucide/svelte";
	import { EmptySlate } from "$components/empty_slate";
	import { resolve } from "$app/paths";

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

	<div>
		{#if orgs.length === 0}
			<EmptySlate caution={true} title="No organisations">
				<a
					href={resolve("/organisations/create")}
					class="text-xs text-muted-foreground underline"
				>
					Add a new organisation to get started
				</a>
			</EmptySlate>
		{:else}
			<Grid>
				{#each orgs as org (org.id)}
					<OrgCard {org} />
				{/each}
			</Grid>
		{/if}
	</div>
</div>
