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
	import { NavItem, NavHeader } from "$components/nav/index.js";
	import { orgSvc } from "$lib/watchtower";
	import { Button } from "$components/ui/button";
	import { Separator } from "$components/ui/separator";
	import { settingsSvc } from "$lib/settings";
	import { BaseTooltip } from "$components/base_tooltip/index.js";

	let { children }: LayoutProps = $props();

	const organisation = $derived(orgSvc.defaultOrg);
	const allOrgs = $derived(orgSvc.organisations);

	let expand = $state(settingsSvc.sidebarExpanded);
	let expandedStyle = $derived(expand ? "w-42" : "w-14");
	let expandIconStyle = $derived(expand ? "justify-end" : "justify-center");

	function toggleExpand(e: Event) {
		e.preventDefault();
		expand = !expand;
		settingsSvc.setSidebarExpanded(expand);
	}
</script>

<div class="flex h-screen">
	<aside
		class={cn(
			"flex h-full flex-col bg-muted p-2 shadow-2xl transition-all duration-300 ease-in-out",
			expandedStyle
		)}
	>
		<div class="flex flex-1 flex-col gap-2">
			{#if organisation?.id}
				<NavHeader {expand} currentOrg={organisation} {allOrgs} />
				<Separator class="mb-2" />
			{/if}
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
		<Separator class="my-2" />
		<div class={cn("flex", expandIconStyle)}>
			<BaseTooltip>
				{#snippet trigger()}
					<Button onclick={toggleExpand} size="sm" variant="ghost">
						{#if expand}
							<PanelLeftClose size={24} />
						{:else}
							<PanelLeftOpen size={24} />
						{/if}
					</Button>
				{/snippet}

				{#if expand}
					<span>Collapse sidebar</span>
				{:else}
					<span>Expand sidebar</span>
				{/if}
			</BaseTooltip>
		</div>
	</aside>
	<main class="flex-1 overflow-auto">
		{@render children?.()}
	</main>
</div>
