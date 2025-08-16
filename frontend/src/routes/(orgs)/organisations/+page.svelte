<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid/index.js";
	import type { PageProps } from "./$types";
	import { OrgCard } from "$components/org_card";

	let { data }: PageProps = $props();

	const orgs = $derived(data.orgs);
</script>

<div class="w-full p-2">
	<PageTitle
		backAction={async () => {
			await goto("/");
		}}
		title="Organisations"
		subtitle="Manage your organisations"
	/>

	<Grid>
		{#each orgs as org (org.id)}
			<OrgCard {org} />
		{/each}
	</Grid>
</div>
