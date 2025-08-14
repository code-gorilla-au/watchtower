<script lang="ts">
	import { type LayoutProps } from "./$types";

	let { children }: LayoutProps = $props();
</script>

<div class="flex">
	<aside class="flex min-h-screen min-w-24 flex-col bg-muted p-2 shadow-sm">
		<div class="flex-1">
			<p>products</p>
			<p>organisations</p>
		</div>
		<p>settings</p>
	</aside>
	{@render children?.()}
</div>
