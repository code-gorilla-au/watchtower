<script lang="ts">
	import { watchtower } from "$lib/wailsjs/go/models";
	import { formatDate } from "$design/formats";
	import { Card, CardContent, CardTitle } from "$components/ui/card";
	import { Badge } from "$components/ui/badge";
	import { Trash } from "@lucide/svelte";
	import { Button } from "$components/ui/button";

	import { goto } from "$app/navigation";
	import { orgSvc } from "$lib/watchtower";

	type Props = {
		org: watchtower.OrganisationDTO;
	};

	let { org }: Props = $props();

	async function syncProduct(id: number) {
		await goto(`/organisations/${id}/edit`);
	}

	async function deleteOrg(id: number) {
		await orgSvc.delete(id);
		await goto("/dashboard");
	}
</script>

<a href={`/products/${org.id}/details`}>
	<Card class="w-full cursor-pointer hover:bg-muted/30">
		<CardTitle class="flex items-center justify-between px-2">
			<span>{org.friendly_name}</span>
			<div>
				<Button
					onclick={async (e: Event) => {
						e.preventDefault();
						await deleteOrg(org.id);
					}}
					size="icon"
					variant="ghost"
				>
					<Trash />
				</Button>
			</div>
		</CardTitle>
		<CardContent>
			<div class="mb-2 flex justify-between text-sm">
				<p class="text-muted-foreground">Last updated:</p>
				<p>{formatDate(org.updated_at)}</p>
			</div>
			{#if org.default_org}
				<Badge>current</Badge>
			{/if}
			<Badge>{org.namespace}</Badge>
		</CardContent>
	</Card>
</a>
