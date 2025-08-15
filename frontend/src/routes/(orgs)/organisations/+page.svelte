<script lang="ts">
	import { PageTitle } from "$components/page_title";
	import { goto } from "$app/navigation";
</script>

<div class="w-full p-2">
	<PageTitle
		backAction={async () => {
			await goto("/");
		}}
		title="Organisations"
		subtitle="Manage your organisations"
	/>
</div>
