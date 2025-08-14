<script lang="ts">
	import { type LayoutProps } from "./$types";

	let { children }: LayoutProps = $props();
</script>

<div class="flex">
	<aside class="flex min-h-screen min-w-24 flex-col bg-muted p-2 shadow-sm">
		<div class="flex flex-1 flex-col gap-2">
			<a href="/products">products</a>
			<a href="/organisations">organisations</a>
		</div>
		<p>settings</p>
	</aside>
	{@render children?.()}
</div>
