<script lang="ts">
	import "../app.css";
	import favicon from "$lib/assets/favicon.svg";
	import { settingsSvc } from "$lib/settings";

	let { children } = $props();
	settingsSvc.initTheme();
</script>

{@render children?.()}

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>
