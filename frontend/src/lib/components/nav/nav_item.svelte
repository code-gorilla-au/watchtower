<script lang="ts">
	import type { Snippet } from "svelte";
	import { page } from "$app/state";
	import { cn } from "$lib/utils";
	import { fade } from "svelte/transition";

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

	let transitionConfig = $derived.by(() => {
		if (expand) {
			return { duration: 150, delay: 100 };
		}

		return { duration: 50, delay: 0 };
	});
</script>

<a
	href={to}
	class={cn(
		"inline-flex gap-2 rounded-md p-2 transition duration-300 ease-in-out",
		className,
		activeLink
	)}
>
	<div>
		{@render icon()}
	</div>
	{#if expand}
		<span transition:fade={transitionConfig} class="w-full">{label}</span>
	{/if}
</a>
