<script lang="ts">
	import { resolve } from "$app/paths";
	import * as DropdownMenu from "$components/ui/dropdown-menu";
	import { fade } from "svelte/transition";
	import { transitionConfig } from "$components/nav/transitions";
	import { organisations } from "$lib/wailsjs/go/models";
	import { orgSvc } from "$lib/watchtower";
	import { goto, invalidateAll } from "$app/navigation";
	type Props = {
		expand?: boolean;
		currentOrg: organisations.OrganisationDTO;
		allOrgs: organisations.OrganisationDTO[];
	};

	let { expand, currentOrg, allOrgs }: Props = $props();

	let orgName = $derived(currentOrg.friendly_name);

	let orgNameInitial = $derived(orgName.charAt(0));
	let trimmedOrgName = $derived.by(() => {
		if (orgName.length <= 12) {
			return orgName;
		}
		return `${orgName.slice(0, 12)}...`;
	});

	let transition = $derived.by(() => {
		if (expand) {
			return transitionConfig.expand;
		}

		return transitionConfig.contract;
	});

	let otherOrgs = $derived.by(() => {
		return allOrgs.filter((o) => o.id !== currentOrg.id);
	});

	async function switchOrganisation(orgId: number) {
		await orgSvc.setDefault(orgId);
		await invalidateAll();
		await goto(resolve("/dashboard"));
	}
</script>

<div class="flex items-center py-2">
	<DropdownMenu.Root>
		<DropdownMenu.Trigger class="flex items-center gap-2">
			{#snippet child({ props })}
				<button {...props} class="inline-flex flex-row items-center gap-2 overflow-auto">
					<span
						class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-secondary p-2 hover:bg-accent hover:text-accent-foreground"
					>
						<span class="capitalize">
							{orgNameInitial}
						</span>
					</span>
					{#if expand}
						<span
							transition:fade={transition}
							class="shrink-0 overflow-hidden font-bold"
						>
							{trimmedOrgName}
						</span>
					{/if}
				</button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content class="w-56" side="bottom" align="start">
			<DropdownMenu.Label>Active Organisation</DropdownMenu.Label>
			<DropdownMenu.Item
				onclick={async () => {
					await switchOrganisation(currentOrg.id);
				}}
			>
				<span class="capitalize">{currentOrg.friendly_name}</span>
			</DropdownMenu.Item>
			{#if otherOrgs.length > 0}
				<DropdownMenu.Separator />
				<DropdownMenu.Label>Other Organisations</DropdownMenu.Label>
				{#each otherOrgs as org (org.id)}
					<DropdownMenu.Item
						onclick={async () => {
							await switchOrganisation(org.id);
						}}
					>
						<a class="w-full" href={resolve(`/organisations/${org.id}`)}>
							<span class="capitalize">{org.friendly_name}</span>
						</a>
					</DropdownMenu.Item>
				{/each}
			{/if}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>
