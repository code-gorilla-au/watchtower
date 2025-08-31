<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Dialog from "$lib/components/ui/dialog";
	import type { Snippet } from "svelte";

	type Props = {
		children: Snippet;
		onConfirm: () => void;
		orgName: string;
	};

	let { onConfirm, orgName, children }: Props = $props();
</script>

<Dialog.Root>
	<Dialog.Trigger>
		{@render children()}
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Confirm Delete organisation: {orgName}?</Dialog.Title>
			<Dialog.Description>
				This will permanently delete your organisation and remove your data, including all
				products, repositories, pull requests, and security reports.
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4 py-4">
			<Button
				onclick={(e: Event) => {
					e.preventDefault();
					onConfirm();
				}}
				variant="destructive">Confirm</Button
			>
		</div>
	</Dialog.Content>
</Dialog.Root>
