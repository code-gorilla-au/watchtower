<script lang="ts">
	import { Button } from "$components/ui/button/index.js";
	import { BaseInput } from "$components/base_input/index.js";
	import type { PageProps } from "./$types";
	import { productSvc } from "$lib/watchtower";
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";

	let { data }: PageProps = $props();

	let organisation = $derived(data.organisation);

	type FormData = {
		name: string;
		tags: string;
		error?: string;
	};

	const form = $state<FormData>({
		name: "",
		tags: "",
		error: undefined
	});

	function onSubmit(e: Event) {
		try {
			e.preventDefault();
			if (!organisation) {
				form.error = "Organisation not found";
				return;
			}

			productSvc.create(form.name, organisation?.id, form.tags);
		} catch (e) {
			form.error = e.message;
		}
	}
</script>

<div class="w-full p-2">
	<PageTitle
		title="Add product"
		backAction={async () => {
			await goto("/");
		}}
		subtitle="Add a product to an organisation {organisation?.friendly_name}"
	/>
	<div class="mb-10"></div>
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
			description="Comma separated list of tags"
			bind:value={form.tags}
		/>
		{#if form.error}
			<div class="rounded-md border border-destructive p-4">
				<h3>Error submitting form</h3>
				<p class="text-destructive">{form.error}</p>
			</div>
		{/if}
		<div class="my-10 flex w-full justify-end">
			<Button type="submit">Add product</Button>
		</div>
	</form>
</div>
