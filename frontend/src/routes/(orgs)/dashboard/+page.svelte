<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";
	import { PRGrid, ProductsGrid } from "$components/products";
	import { onDestroy, onMount } from "svelte";
	import { TIME_TWO_MINUTES } from "$lib/watchtower/types";
	import { invalidateAll } from "$app/navigation";
	import { SecurityGrid } from "$components/products/index.js";

	let { data }: PageProps = $props();
	let org = $derived(data.organisation);
	let products = $derived(data.products);
	let prs = $derived(data.prs);
	let securities = $derived(data.securities);

	let intervalPoll: number;
	onMount(() => {
		intervalPoll = setInterval(async () => {
			await invalidateAll();
		}, TIME_TWO_MINUTES);
	});
	onDestroy(() => {
		clearInterval(intervalPoll);
	});
</script>

<div class="page-container">
	<PageTitle title="Dashboard" subtitle={org?.friendly_name} />

	<div class="my-4">
		<h3 class="text-xl text-muted-foreground">Security Vulnerability</h3>
	</div>
	<SecurityGrid {securities} />

	<div class="my-4">
		<h3 class="text-xl text-muted-foreground">Pull Requests</h3>
	</div>
	<PRGrid {prs} />

	<div class="my-4">
		<h3 class="text-xl text-muted-foreground">Products</h3>
	</div>
	<ProductsGrid {products} />
</div>
