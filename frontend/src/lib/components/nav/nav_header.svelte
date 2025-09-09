<script lang="ts">
	import { resolve } from "$app/paths";
	import * as DropdownMenu from "$components/ui/dropdown-menu";
	import { fade } from "svelte/transition";
	import { transitionConfig } from "$components/nav/transitions";
	import { watchtower } from "$lib/wailsjs/go/models";
	type Props = {
		expand?: boolean;
		currentOrg: watchtower.OrganisationDTO;
		allOrgs: watchtower.OrganisationDTO[];
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
</script>

<div class="flex items-center py-2">
	<DropdownMenu.Root>
		<DropdownMenu.Trigger class="flex items-center gap-2">
			{#snippet child({ props })}
				<button {...props} class="inline-flex items-center gap-2">
					<span
						class="flex h-8 w-8 items-center justify-center rounded-full bg-secondary p-2 hover:bg-accent hover:text-accent-foreground"
					>
						<span class="capitalize">
							{orgNameInitial}
						</span>
					</span>
					{#if expand}
						<span transition:fade={transition} class="font-bold">
							{trimmedOrgName}
						</span>
					{/if}
				</button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content class="w-56" side="bottom" align="start">
			<DropdownMenu.Label>Active Organisation</DropdownMenu.Label>
			<DropdownMenu.Item>
				<a class="w-full" href={resolve(`/organisations/${currentOrg.id}`)}>
					<span class="capitalize">{currentOrg.friendly_name}</span>
				</a>
			</DropdownMenu.Item>
			<DropdownMenu.Separator />
			<DropdownMenu.Label>Other Organisations</DropdownMenu.Label>
			{#each otherOrgs as org (org.id)}
				<DropdownMenu.Item>
					<a class="w-full" href={resolve(`/organisations/${org.id}`)}>
						<span class="capitalize">{org.friendly_name}</span>
					</a>
				</DropdownMenu.Item>
			{/each}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>
