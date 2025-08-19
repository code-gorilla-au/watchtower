<script lang="ts">
	import type { Snippet } from "svelte";
	import { cn } from "$lib/utils";

	type Props = {
		class?: string;
		expand: boolean;
		to: string;
		label: string;
		icon: Snippet;
	};
	let { expand, to, label, icon, class: className }: Props = $props();
</script>

<a href={to} class={cn("mb-2 flex items-center justify-center gap-2", className)}>
	<span>
		{@render icon()}
	</span>
	{#if expand}
		<span class="w-full">{label}</span>
	{/if}
</a>
