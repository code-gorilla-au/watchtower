<script lang="ts">
	import type { Snippet } from "svelte";
	import type { HTMLButtonAttributes } from "svelte/elements";

	type Props = HTMLButtonAttributes & {
		class?: string;
		children?: Snippet;
	};

	let { class: className, children }: Props = $props();
</script>

<button class={className}>{children?.()}</button>
