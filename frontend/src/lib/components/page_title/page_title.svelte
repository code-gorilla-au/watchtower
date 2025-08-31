<script lang="ts">
	import { ChevronLeft } from "@lucide/svelte";
	import { Button } from "$components/ui/button";
	import type { Snippet } from "svelte";
	import { cn } from "$lib/utils";

	export type Props = {
		backAction?: () => void;
		title: string;
		subtitle?: string;
		class?: string;
		children?: Snippet;
	};
	let { backAction, title, subtitle, children, class: className }: Props = $props();
</script>

<div class={cn("mb-10 flex items-center justify-between", className)}>
	<div class="flex items-center gap-2">
		{#if backAction}
			<Button onclick={backAction} size="icon" variant="ghost">
				<ChevronLeft />
			</Button>
		{/if}
		<div>
			<h1 class="heading-1 mb-2">{title}</h1>
			{#if subtitle}
				<p class="text-sm">{subtitle}</p>
			{/if}
		</div>
	</div>
	<div>
		{@render children?.()}
	</div>
</div>
