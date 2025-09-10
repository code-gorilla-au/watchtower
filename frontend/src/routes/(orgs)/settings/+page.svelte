<script lang="ts">
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
	import { settingsSvc } from "$lib/settings";
	import { Label } from "$components/ui/label";
	import { Switch } from "$components/ui/switch";
	import { Separator } from "$components/ui/separator";
	import { Button } from "$components/ui/button";
	import { Trash } from "@lucide/svelte";
	import { orgSvc } from "$lib/watchtower";
	import { resolve } from "$app/paths";

	type FormState = {
		darkMode: boolean;
	};

	const formState = $state<FormState>({
		darkMode: settingsSvc.theme === "dark"
	});

	$effect(() => {
		settingsSvc.setTheme(formState.darkMode ? "dark" : "light");
	});

	async function deleteAllData(e: Event) {
		e.preventDefault();

		await orgSvc.deleteAll();
		await goto(resolve("/register/organisation"));
	}
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto(resolve("/"));
		}}
		title="Settings"
		subtitle="User general settings"
	/>
	<form class="mx-auto flex max-w-lg flex-col gap-4">
		<h3 class="heading-2">User preferences</h3>
		<div class="flex w-full items-center justify-between">
			<Label for="darkMode">Switch to dark mode</Label>
			<Switch id="darkMode" bind:checked={formState.darkMode} />
		</div>
		<Separator class="my-5" />
		<h3 class="heading-2">System</h3>
		<div class="flex w-full items-center justify-between">
			<p>App version:</p>
			<p>{settingsSvc.version}</p>
		</div>
		<div class="flex w-full items-center justify-between">
			<p>Local database location:</p>
			<p>{settingsSvc.version}</p>
		</div>
		<Separator class="my-5" />
		<h3 class="heading-2">Danger zone</h3>
		<div class="flex w-full items-center justify-between">
			<div>
				<p>Delete all data</p>
				<p class="text-sm text-muted-foreground">This action cannot be undone.</p>
			</div>
			<Button onclick={deleteAllData} variant="destructive" size="icon">
				<Trash />
			</Button>
		</div>
	</form>
</div>
