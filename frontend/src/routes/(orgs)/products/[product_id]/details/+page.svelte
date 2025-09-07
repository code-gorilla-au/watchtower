<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { RefreshCw, Pencil } from "@lucide/svelte";
	import { type PageProps } from "./$types";
	import { goto } from "$app/navigation";
	import { Grid } from "$components/grid";
	import { RepoCard } from "$components/products";
	import { Button } from "$components/ui/button";
	import { PRGrid, SecurityGrid } from "$components/products/index.js";
	import * as Accordion from "$components/ui/accordion";
	import { TagsFilter } from "$lib/hooks/filters.svelte";
	import { GridHeader } from "$components/grid/index.js";

	let { data }: PageProps = $props();

	const prsFilter = new TagsFilter(data.prs, "repository_name");

	let product = $derived(data.product);
	let repos = $derived(data.repos);
	let prs = $derived(data.prs);
	let securities = $derived(data.securities);

	async function syncProduct(e: Event) {
		e.preventDefault();

		await goto(`/products/${product.id}/sync`);
	}

	async function editProduct(e: Event) {
		e.preventDefault();
		await goto(`/products/${product.id}/edit`);
	}
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(`/products`);
		}}
		title="Product Details"
		subtitle="Details for product {product.name}"
	>
		<Button onclick={syncProduct} variant="ghost" size="icon">
			<RefreshCw />
		</Button>
		<Button onclick={editProduct} variant="ghost" size="icon">
			<Pencil />
		</Button>
	</PageTitle>

	<Accordion.Root type="multiple" value={["sec", "prs", "repos"]}>
		<Accordion.Item value="sec">
			<Accordion.Trigger>
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
			<Accordion.Trigger>
				<GridHeader data={prs} tagField="repository_name" title="Pull Requests" />
			</Accordion.Trigger>
			<Accordion.Content class="mb-5">
				<PRGrid prs={prsFilter.data} />
			</Accordion.Content>
		</Accordion.Item>
		<Accordion.Item value="repos">
			<Accordion.Trigger>
				<h2 class="text-xl text-muted-foreground">Repositories ({repos.length})</h2>
			</Accordion.Trigger>
			<Accordion.Content class="mb-5">
				<Grid>
					{#each repos as repo (repo.id)}
						<RepoCard {repo} />
					{/each}
				</Grid>
			</Accordion.Content>
		</Accordion.Item>
	</Accordion.Root>
</div>
