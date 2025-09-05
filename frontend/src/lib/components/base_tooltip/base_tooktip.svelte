<script lang="ts">
	import * as Tooltip from "$lib/components/ui/tooltip";
	import type { Snippet } from "svelte";

	type Props = {
		trigger: Snippet;
		children?: Snippet;
	};

	let { children, trigger }: Props = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger>{@render trigger()}</Tooltip.Trigger>
		<Tooltip.Content>
			{@render children?.()}
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
