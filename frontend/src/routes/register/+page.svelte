<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { OpenExternalURL } from "$lib/wailsjs/go/main/App";
	import { Button } from "$components/ui/button";

	async function routeToPAT() {
		await OpenExternalURL(
			"https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens"
		);
	}

	async function routeToTopics() {
		await OpenExternalURL(
			"https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-topics"
		);
	}
</script>

<div class="page-container">
	<PageTitle title="Register" subtitle="Getting started" />
	<div class="mx-auto flex max-w-lg flex-col text-muted-foreground">
		<h2 class="heading-2 mb-2">What is needed to begin</h2>
		<p class="mb-8">
			A <button class="cursor-pointer underline" onclick={routeToPAT}
				>Personal Access Token (PAT)</button
			>
			is required to be registered when onboarding an organisation. Watchtower requires
			<b> readonly </b>
			access to repos, no data leaves the device.
		</p>

		<h2 class="heading-2 mb-2">Tagging repositories</h2>
		<p>
			Organising repositories, pull requests and security advisories within watchtower is done
			using <button class="cursor-pointer underline" onclick={routeToTopics}>
				Github topics.
			</button>
			Tag all repositories you wish to monitor with the <b>watchtower</b> topic.
		</p>
		<div class="my-8 flex justify-end">
			<a href="/register/organisation" class="cursor-pointer underline"
				>Onboard my first organisation</a
			>
		</div>
	</div>
</div>
