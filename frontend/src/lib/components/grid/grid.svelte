<script lang="ts">
	import { cn } from "$lib/utils";
	import type { Snippet } from "svelte";

	let { children, class: className }: { children: Snippet; class?: string } = $props();
</script>

<div class={cn("grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-4", className)}>
	{@render children()}
</div>
