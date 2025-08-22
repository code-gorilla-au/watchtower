<script lang="ts">
	import { Button } from "$components/ui/button/index.js";
	import { goto } from "$app/navigation";

	async function backToDashboard(e: Event) {
		e.preventDefault();

		await goto("/");
	}
</script>

<h1 class="text-4xl">Error</h1>
<Button onclick={backToDashboard}>Back to dashboard</Button>
