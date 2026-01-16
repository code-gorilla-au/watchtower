<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";
	import { ProductsGrid, AccordionItemSecurity, AccordionItemPrs } from "$components/products";
	import { onDestroy, onMount } from "svelte";
	import { TIME_TWO_MINUTES } from "$lib/watchtower/types";
	import { invalidateAll } from "$app/navigation";
	import { TimeSince } from "$lib/hooks/time.svelte";
	import { SvelteDate } from "svelte/reactivity";
	import * as Accordion from "$components/ui/accordion";
	import { Button } from "$components/ui/button";
	import { Search } from "@lucide/svelte";
	import { SearchBar } from "$components/searchbar";

	let { data }: PageProps = $props();
	let org = $derived(data.organisation);
	let products = $derived(data.products);
	let prs = $derived(data.prs);
	let securities = $derived(data.securities);

	let intervalPoll: NodeJS.Timeout | undefined;

	const timeSince = new TimeSince(new SvelteDate());

	let searchBarOpen = $state(false);

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
	<SearchBar bind:open={searchBarOpen} {securities} {prs} />
	<PageTitle title="Dashboard" subtitle={org?.description || org?.friendly_name}>
		<div class="flex items-center gap-2">
			<p class="text-xs text-muted-foreground">Last sync: {timeSince.date}</p>
			<Button
				onclick={(e) => {
					e.preventDefault();
					searchBarOpen = !searchBarOpen;
				}}
				class="hover:text-accent"
				variant="ghost"
			>
				<Search />
			</Button>
		</div>
	</PageTitle>
	<Accordion.Root type="multiple" value={["security", "prs", "products"]}>
		<AccordionItemSecurity {securities} />
		<AccordionItemPrs {prs} />
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
