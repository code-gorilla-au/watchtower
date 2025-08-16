<script lang="ts">
	import type { PageProps } from "./$types";
	import { PageTitle } from "$components/page_title";

	let { data }: PageProps = $props();
	let org = $derived(data.organisation);
</script>

<div class="p-2">
	<PageTitle title="Dashboard" subtitle={org?.friendly_name} />
	<p>{JSON.stringify(data)}</p>
</div>
