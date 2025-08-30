<script lang="ts">
	type Props = {
		expand?: boolean;
		orgName: string;
		orgId: string;
	};

	let { expand, orgName, orgId }: Props = $props();

	let orgNameInitial = $derived(orgName.charAt(0));
</script>

<div class="flex items-center py-2">
	<a href="/organisations/{orgId}" class="flex w-full items-center">
		{#if expand}
			<span class="p-2 font-bold">
				{orgName}
			</span>
		{:else}
			<span class="w-full rounded-full bg-secondary py-2 text-center">
				{orgNameInitial}
			</span>
		{/if}
	</a>
</div>
