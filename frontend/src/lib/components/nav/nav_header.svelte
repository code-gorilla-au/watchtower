<script lang="ts">
	type Props = {
		expand?: boolean;
		orgName: string;
	};

	let { expand, orgName }: Props = $props();

	let isExpanded = $derived(expand);
	let renderOrgName = $derived.by(() => {
		if (isExpanded) {
			return orgName;
		}
		return orgName.charAt(0);
	});
</script>

<div class="flex items-center py-2">
	{#if expand}
		<span class="py-2 font-bold">
			{renderOrgName}
		</span>
	{:else}
		<span class="w-full rounded-full bg-secondary py-2 text-center">
			{renderOrgName}
		</span>
	{/if}
</div>
