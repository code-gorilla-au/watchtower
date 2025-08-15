<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";

	let { data }: PageProps = $props();

	let product = $derived(data.product);

	const syncState = $state({
		loading: false
	});
</script>

<div class="w-full p-2">
	<PageTitle title="Sync Product {product.name}" subtitle="Sync pull request from GitHub" />
</div>
