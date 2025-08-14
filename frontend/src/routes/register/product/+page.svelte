<script lang="ts">
	import { Button } from "$components/ui/button/index.js";
	import { BaseInput } from "$components/base_input/index.js";
	import type { PageProps } from "./$types";

	let { data }: PageProps = $props();

	let organisation = $derived(data.organisation);

	type FormData = {
		name: string;
		tags: string;
	};

	const form = $state<FormData>({
		name: "",
		tags: ""
	});

	function onSubmit(e: Event) {
		e.preventDefault();
	}
</script>

<div class="p-3">
	<h1 class="text-4xl">Register - Product</h1>
	<div class="mb-10">Add a product to an organisation {organisation.friendly_name}</div>
	<form
		method="POST"
		onsubmit={onSubmit}
		class="mx-auto flex max-w-sm flex-col items-center justify-center"
	>
		<BaseInput
			class=""
			id="name"
			label="Product name"
			description="Product friendly name"
			bind:value={form.name}
		/>

		<BaseInput
			id="tags"
			label="Tags"
			description="comma separated list of tags"
			bind:value={form.tags}
		/>

		<div class="my-10 flex w-full justify-end">
			<Button type="submit">Sync product</Button>
		</div>
	</form>
</div>
