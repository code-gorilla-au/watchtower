<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";
	import { PRGrid, ProductsGrid } from "$components/products";
	import { GridHeader } from "$components/grid/index.js";
	import { onDestroy, onMount } from "svelte";
	import { TIME_TWO_MINUTES } from "$lib/watchtower/types";
	import { invalidateAll } from "$app/navigation";
	import { SecurityGrid } from "$components/products/index.js";
	import { TimeSince } from "$lib/hooks/time.svelte";
	import { SvelteDate } from "svelte/reactivity";
	import * as Accordion from "$components/ui/accordion";

	let { data }: PageProps = $props();
	let org = $derived(data.organisation);
	let products = $derived(data.products);
	let prs = $derived(data.prs);
	let securities = $derived(data.securities);

	let intervalPoll: NodeJS.Timeout | undefined;

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
	<PageTitle title="Dashboard" subtitle={org?.description || org?.friendly_name}>
		<p class="text-xs text-muted-foreground">Last sync: {timeSince.date}</p>
	</PageTitle>
	<Accordion.Root type="multiple" value={["security", "prs", "products"]}>
		<Accordion.Item value="security">
			<Accordion.Trigger class="text-left">
				<GridHeader
					data={securities}
					tagField="repository_name"
					title="Security Vulnerabilities"
				/>
			</Accordion.Trigger>
			<Accordion.Content class="mb-5">
				<SecurityGrid {securities} />
			</Accordion.Content>
		</Accordion.Item>
		<Accordion.Item value="prs">
			<Accordion.Trigger class="text-left">
				<GridHeader data={prs} tagField="repository_name" title="Pull Requests" />
			</Accordion.Trigger>
			<Accordion.Content class="mb-5">
				<PRGrid {prs} />
			</Accordion.Content>
		</Accordion.Item>
		<Accordion.Item value="products">
			<Accordion.Trigger class="text-left">
				<h3 class="text-xl text-muted-foreground">Products ({products.length})</h3>
			</Accordion.Trigger>
			<Accordion.Content class="mb-5">
				<ProductsGrid {products} />
			</Accordion.Content>
		</Accordion.Item>
	</Accordion.Root>
</div>
