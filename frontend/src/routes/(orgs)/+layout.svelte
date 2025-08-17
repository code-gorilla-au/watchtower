<script lang="ts">
	import { type LayoutProps } from "./$types";
	import { Settings, Package, Castle, LayoutDashboard } from "@lucide/svelte";
	import { cn } from "$lib/utils";
	import { NavItem } from "$components/nav/index.js";

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
			<NavItem to="/" {expand} label="Dashboard">
				{#snippet icon()}
					<LayoutDashboard size={24} />
				{/snippet}
			</NavItem>
			<NavItem to="/products" {expand} label="Products">
				{#snippet icon()}
					<Package size={24} />
				{/snippet}
			</NavItem>
			<NavItem to="/organisations" {expand} label="Organisations">
				{#snippet icon()}
					<Castle size={24} />
				{/snippet}
			</NavItem>
		</div>
		<NavItem {expand} to="/settings" label="Settings">
			{#snippet icon()}
				<Settings size={24} />
			{/snippet}
		</NavItem>
	</aside>
	{@render children?.()}
</div>
