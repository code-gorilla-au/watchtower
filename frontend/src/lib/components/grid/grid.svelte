<script lang="ts">
	import { cn } from "$lib/utils";
	import type { Snippet } from "svelte";

	let { children, class: className }: { children: Snippet; class?: string } = $props();
</script>

<div class={cn("grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4", className)}>
	{@render children()}
</div>
