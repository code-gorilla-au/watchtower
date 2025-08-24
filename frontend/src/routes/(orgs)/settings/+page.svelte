<script lang="ts">
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { settingsSvc } from "$lib/settings";
	import { Label } from "$components/ui/label";
	import { Switch } from "$components/ui/switch";

	type FormState = {
		darkMode: boolean;
	};

	const formState = $state<FormState>({
		darkMode: settingsSvc.theme === "dark"
	});
	$effect(() => {
		settingsSvc.setTheme(formState.darkMode ? "dark" : "light");
	});
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto("/");
		}}
		title="Settings"
		subtitle="User general settings"
	/>
	<form class="mx-auto max-w-lg">
		<h3 class="mb-10 text-2xl">User preferences</h3>
		<div class="flex w-full items-center justify-between gap-4">
			<Label for="darkMode">Switch to dark mode</Label>
			<Switch id="darkMode" bind:checked={formState.darkMode} />
		</div>
	</form>
</div>
