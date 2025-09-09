<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";
	import { onMount } from "svelte";
	import { productSvc } from "$lib/watchtower";
	import { goto } from "$app/navigation";
	import { EmptySlate } from "$components/empty_slate/index.js";
	import { LoaderSquare } from "$components/loaders";
	import { resolve } from "$app/paths";

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

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(resolve(`/products`));
		}}
		title="Sync Product {product.name}"
		subtitle="Sync pull request from GitHub"
	/>

	{#if syncState.loading}
		<div class="flex flex-col items-center justify-center">
			<LoaderSquare />
			<h2 class="heading-2">Syncing product please wait</h2>
		</div>
	{:else if syncState.error}
		<div>
			<h3>Error</h3>
			<p>{syncState.error}</p>
		</div>
	{:else}
		<div class="flex flex-col items-center justify-center">
			<EmptySlate class="w-full" title="Product sync complete">
				<a class="underline" href={resolve(`/products/${product.id}/details`)}>
					View product details
				</a>
			</EmptySlate>
		</div>
	{/if}
</div>
