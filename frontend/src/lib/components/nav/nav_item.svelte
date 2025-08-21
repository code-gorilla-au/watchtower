<script lang="ts">
	import type { Snippet } from "svelte";
	import { page } from "$app/state";
	import { cn } from "$lib/utils";

	type Props = {
		class?: string;
		expand: boolean;
		to: string;
		label: string;
		icon: Snippet;
	};

	let { expand, to, label, icon, class: className }: Props = $props();

	let currentActive = $derived(page.url.pathname.startsWith(to));
	let activeLink = $derived(
		currentActive ? "bg-accent text-accent-foreground" : "hover:bg-secondary"
	);
</script>

<a
	href={to}
	class={cn(
		"flex items-center justify-center gap-2 rounded-md p-2 transition duration-300 ease-in-out",
		className,
		activeLink
	)}
>
	<span>
		{@render icon()}
	</span>
	{#if expand}
		<span class="w-full">{label}</span>
	{/if}
</a>
