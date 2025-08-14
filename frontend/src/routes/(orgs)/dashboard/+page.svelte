<script lang="ts">
	import type { PageProps } from "./$types";

	let { data }: PageProps = $props();
</script>

<div>
	<h1>Dashboard</h1>
	<span>{data?.organisation?.friendly_name}</span>
</div>

<p>{JSON.stringify(data)}</p>
