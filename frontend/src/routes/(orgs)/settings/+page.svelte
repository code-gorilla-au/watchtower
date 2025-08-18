<script lang="ts">
	import { PageTitle } from "$components/page_title/index.js";
	import { goto } from "$app/navigation";
</script>

<div class="page-container">
	<PageTitle
		backAction={async () => {
			await goto("/");
		}}
		title="Settings"
		subtitle="Organisation settings"
	/>
</div>
