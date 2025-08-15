<script lang="ts">
	import { PageTitle } from "$components/page_title";

	const syncState = $state({
		loading: false
	});
</script>

<div class="w-full p-2">
	<PageTitle title="Sync Product" subtitle="Sync pull request from GitHub" />
</div>
