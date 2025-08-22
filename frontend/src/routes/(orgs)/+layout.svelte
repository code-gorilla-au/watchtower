<script lang="ts">
	import { type LayoutProps } from "./$types";
	import {
		Settings,
		Package,
		Castle,
		LayoutDashboard,
		PanelLeftClose,
		PanelLeftOpen
	} from "@lucide/svelte";
	import { cn } from "$lib/utils";
	import { NavItem } from "$components/nav/index.js";
	import { orgSvc } from "$lib/watchtower";
	import { Button } from "$components/ui/button";
	import { Separator } from "$components/ui/separator";
	import { settingsSvc } from "$lib/settings";

	let { children }: LayoutProps = $props();

	const organisation = $derived(orgSvc.defaultOrg);

	let expand = $state(settingsSvc.sidebarExpanded);
	let expandedStyle = $derived(expand ? "min-w-40" : "w-14");

	function toggleExpand(e: Event) {
		e.preventDefault();
		expand = !expand;
		settingsSvc.setSidebarExpanded(expand);
	}
</script>

<div class="flex h-screen">
	<aside class={cn("flex h-full max-w-40 flex-col bg-muted p-2 shadow-2xl ", expandedStyle)}>
		<div class="flex flex-1 flex-col gap-2">
			<div class="flex items-center justify-between py-2">
				{#if expand}
					<span class="font-bold">
						{organisation?.friendly_name}
					</span>
				{/if}

				<Button onclick={toggleExpand} size="sm" variant="ghost">
					{#if expand}
						<PanelLeftClose />
					{:else}
						<PanelLeftOpen />
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
