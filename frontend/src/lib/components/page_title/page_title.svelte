<script lang="ts">
	import { ChevronLeft } from "@lucide/svelte";
	import { Button } from "$components/ui/button";
	import type { Snippet } from "svelte";

	export type Props = {
		backAction?: () => void;
		title: string;
		subtitle?: string;
		children?: Snippet;
	};
	let { backAction, title, subtitle, children }: Props = $props();
</script>

<div class="flex justify-between">
	<div class="flex items-center gap-2">
		{#if backAction}
			<Button onclick={backAction} size="icon" variant="ghost">
				<ChevronLeft />
			</Button>
		{/if}
		<div>
			<h1 class="text-2xl">{title}</h1>
			{#if subtitle}
				<p class="text-xs">{subtitle}</p>
			{/if}
		</div>
	</div>
	<div>
		{@render children?.()}
	</div>
</div>
