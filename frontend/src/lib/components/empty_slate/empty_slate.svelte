<script lang="ts">
	import { CircleCheckBig, ShieldAlert } from "@lucide/svelte";
	import type { Snippet } from "svelte";

	type EmptySlateProps = {
		title: string;
		description?: string;
		caution?: boolean;
		children?: Snippet;
	};

	const { title, description, caution = false, children }: EmptySlateProps = $props();
</script>

<div class=" flex flex-col items-center justify-center rounded-md bg-muted p-10 text-sm">
	{#if caution}
		<ShieldAlert class="text-muted-foreground" size={62} />
	{:else}
		<CircleCheckBig class="text-green-700" size={62} />
	{/if}

	<h3 class="my-2 text-lg">{title}</h3>
	{#if description}
		<p class=" text-muted-foreground">{description}</p>
	{/if}
	{@render children?.()}
</div>
