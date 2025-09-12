<script lang="ts">
	import "../app.css";
	import favicon from "$lib/assets/favicon.svg";
	import { Toaster } from "$components/ui/sonner";
	import { settingsSvc } from "$lib/settings";
	import { onDestroy } from "svelte";
	import { orgSvc } from "$lib/watchtower";

	let { children } = $props();
	settingsSvc.initTheme();

	onDestroy(() => {
		orgSvc.unmount();
	});
</script>

<Toaster />

{@render children?.()}

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>
