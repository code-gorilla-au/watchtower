<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";
	import { onMount } from "svelte";
	import { productSvc } from "$lib/watchtower";

	let { data }: PageProps = $props();

	let product = $derived(data.product);

	const syncState = $state({
		loading: false,
		error: ""
	});

	onMount(async () => {
		try {
			syncState.loading = true;
			await productSvc.syncProduct(product.id);
		} catch (e) {
			const err = e as Error;
			syncState.error = err.message;
		} finally {
			syncState.loading = false;
		}
	});
</script>

<div class="w-full p-2">
	<PageTitle title="Sync Product {product.name}" subtitle="Sync pull request from GitHub" />

	{#if syncState.loading}
		<div>
			<h3>Loading...</h3>
		</div>
	{/if}
</div>
