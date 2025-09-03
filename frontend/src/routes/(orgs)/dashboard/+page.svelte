<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";
	import { PRGrid, ProductsGrid } from "$components/products";
	import { onDestroy, onMount } from "svelte";
	import { TIME_TWO_MINUTES } from "$lib/watchtower/types";
	import { invalidateAll } from "$app/navigation";
	import { SecurityGrid } from "$components/products/index.js";
	import { TimeSince } from "$lib/hooks/time.svelte";
	import { SvelteDate } from "svelte/reactivity";

	let { data }: PageProps = $props();
	let org = $derived(data.organisation);
	let products = $derived(data.products);
	let prs = $derived(data.prs);
	let securities = $derived(data.securities);

	let intervalPoll: number;

	let timeSince = new TimeSince(new SvelteDate());

	onMount(() => {
		timeSince.start();

		intervalPoll = setInterval(async () => {
			await invalidateAll();

			timeSince.setDate(new SvelteDate());
		}, TIME_TWO_MINUTES);
	});

	onDestroy(() => {
		clearInterval(intervalPoll);

		timeSince.stop();
	});
</script>

<div class="page-container">
	<PageTitle title="Dashboard" subtitle={org?.friendly_name}>
		<p class="text-xs text-muted-foreground">Last sync: {timeSince.date}</p>
	</PageTitle>
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
