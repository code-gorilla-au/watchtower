<script lang="ts">
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import type { HTMLInputAttributes, HTMLInputTypeAttribute } from "svelte/elements";
	import type { WithElementRef } from "$lib/utils";

	type InputType = Exclude<HTMLInputTypeAttribute, "file">;

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, "type"> &
			(
				| { type: "file"; files?: FileList; label?: string; description?: string }
				| { type?: InputType; files?: undefined; label?: string; description?: string }
			)
	>;

	let { id, label, description, ...restProps }: Props = $props();
</script>

<div class="flex w-full max-w-sm flex-col gap-1.5">
	{#if label}
		<Label for={id}>{label}</Label>
	{/if}
	<Input {id} {...restProps} />
	{#if description}
		<p class="text-sm text-muted-foreground">{description}</p>
	{/if}
</div>
