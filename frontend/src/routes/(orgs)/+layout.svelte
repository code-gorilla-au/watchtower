<script lang="ts">
	import { type LayoutProps } from "./$types";
	import {
		Settings,
		Package,
		Castle,
		LayoutDashboard,
		ChevronLeft,
		ChevronRight
	} from "@lucide/svelte";
	import { cn } from "$lib/utils";
	import { NavItem } from "$components/nav/index.js";
	import { orgSvc } from "$lib/watchtower";
	import { Button } from "$components/ui/button";
	import { Separator } from "$components/ui/separator";

	let { children }: LayoutProps = $props();

	const organisation = $derived(orgSvc.defaultOrg);

	let expand = $state(true);
	let expandedStyle = $derived(expand ? "min-w-40" : "w-14");
	function toggleExpand(e: Event) {
		e.preventDefault();
		expand = !expand;
	}
</script>

<div class="flex h-screen">
	<aside
		class={cn(
			"flex h-full max-w-40 flex-col bg-muted p-2 shadow-sm transition-all duration-300 ease-in-out",
			expandedStyle
		)}
	>
		<div class="flex flex-1 flex-col gap-2">
			<div class="py-2">
				{#if expand}
					<span class="font-bold">
						{organisation?.friendly_name}
					</span>
				{/if}

				<Button onclick={toggleExpand} size="sm" variant="ghost" class="ml-auto">
					{#if expand}
						<ChevronLeft />
					{:else}
						<ChevronRight />
					{/if}
				</Button>
			</div>
			<Separator class="mb-2" />
			<NavItem to="/dashboard" {expand} label="Dashboard">
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
	<main class="flex-1 overflow-auto">
		{@render children?.()}
	</main>
</div>
