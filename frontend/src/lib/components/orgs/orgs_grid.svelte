<script lang="ts">
	import { resolve } from "$app/paths";
	import { OrgCard } from "$components/orgs/org_card";
	import { EmptySlate } from "$components/empty_slate";
	import { Grid } from "$components/grid";
	import { organisations } from "$lib/wailsjs/go/models";
	import { orgSvc } from "$lib/watchtower";
	import { goto } from "$app/navigation";
	import { toast } from "svelte-sonner";

	type Props = {
		orgs: organisations.OrganisationDTO[];
	};

	let { orgs }: Props = $props();

	function findOrgById(id: number) {
		return orgs.find((o) => o.id === id);
	}

	async function deleteOrg(id: number) {
		const org = findOrgById(id);

		await orgSvc.delete(id);
		toast.success(`Organisation ${org?.friendly_name} deleted`, {
			position: "top-right"
		});
		await goto(resolve("/dashboard"));
	}
</script>

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
				<OrgCard {org} onDelete={deleteOrg} />
			{/each}
		</Grid>
	{/if}
</div>
