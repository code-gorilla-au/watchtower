<script lang="ts">
	import { type LayoutProps } from "./$types";
	import { Settings, Package } from "@lucide/svelte";
	import { cn } from "$lib/utils";

	let { children }: LayoutProps = $props();
	let expand = $state(true);
	let expandedStyle = $derived(expand ? "min-w-40" : "w-14");
</script>

<div class="flex">
	<aside
		class={cn(
			"flex min-h-screen max-w-40 flex-col bg-muted p-2 shadow-sm transition-all duration-300 ease-in-out",
			expandedStyle
		)}
	>
		<div class="flex flex-1 flex-col gap-2">
			<a href="/dashboard">dashboard</a>
			<a href="/products" class="flex gap-2">
				<Package size={24} />
				{#if expand}
					<span>Products</span>
				{/if}
			</a>
			<a href="/organisations">organisations</a>
		</div>
		<a href="/settings" class="flex justify-center gap-2 pb-2">
			<Settings size={24} />
			{#if expand}
				<span>Settings</span>
			{/if}
		</a>
	</aside>
	{@render children?.()}
</div>
