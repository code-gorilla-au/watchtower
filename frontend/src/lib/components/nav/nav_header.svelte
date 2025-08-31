<script lang="ts">
	import { fade } from "svelte/transition";
	import { transitionConfig } from "$components/nav/transitions";
	type Props = {
		expand?: boolean;
		orgName: string;
		orgId: number;
	};

	let { expand, orgName, orgId }: Props = $props();

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
</script>

<div class="flex items-center py-2">
	<a href="/organisations/{orgId}" class="text-nowrap">
		<div class="inline-flex items-center gap-2">
			<div
				class="flex h-8 w-8 items-center justify-center rounded-full bg-secondary p-2 hover:bg-accent hover:text-accent-foreground"
			>
				<span class="capitalize">
					{orgNameInitial}
				</span>
			</div>
			{#if expand}
				<span transition:fade={transition} class="font-bold">
					{trimmedOrgName}
				</span>
			{/if}
		</div>
	</a>
</div>
