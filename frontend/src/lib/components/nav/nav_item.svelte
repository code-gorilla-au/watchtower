<script lang="ts">
	import type { Snippet } from "svelte";
	import { page } from "$app/state";
	import { cn } from "$lib/utils";
	import { fade } from "svelte/transition";
	import { transitionConfig } from "$components/nav/transitions";

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

	let transition = $derived.by(() => {
		if (expand) {
			return transitionConfig.expand;
		}

		return transitionConfig.contract;
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
		<span transition:fade={transition} class="w-full">{label}</span>
	{/if}
</a>
