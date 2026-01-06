<script lang="ts">
	import { SecurityCard } from "$components/products/index";
	import { EmptySlate } from "$components/empty_slate";
	import { Grid } from "$components/grid";
	import { products } from "$lib/wailsjs/go/models";

	type Props = {
		securities: products.SecurityDTO[];
	};
	let { securities }: Props = $props();
</script>

{#if securities.length === 0}
	<EmptySlate title="No security vulnerabilities"></EmptySlate>
{:else}
	<Grid>
		{#each securities as security (security.id)}
			<SecurityCard {security} />
		{/each}
	</Grid>
{/if}
