<script lang="ts">
	import { organisations } from "$lib/wailsjs/go/models";
	import { formatDate } from "$lib/hooks/formats";
	import { Card, CardContent, CardTitle } from "$components/ui/card";
	import { Badge } from "$components/ui/badge";
	import { Trash } from "@lucide/svelte";
	import { Button } from "$components/ui/button";
	import OrgRemoveDialogue from "$components/orgs/org_remove_dialog.svelte";
	import { CardAction, CardHeader } from "$components/ui/card/index.js";
	import { resolve } from "$app/paths";

	type Props = {
		org: organisations.OrganisationDTO;
		onDelete?: (id: number) => void;
	};

	let { org, onDelete }: Props = $props();

	async function deleteOrg(id: number) {
		onDelete?.(id);
	}
</script>

<a href={resolve(`/organisations/${org.id}`)}>
	<Card class="w-full cursor-pointer hover:bg-muted/30">
		<CardHeader class="flex items-center justify-between">
			<CardTitle>
				<span>{org.friendly_name}</span>
			</CardTitle>
			<CardAction>
				<OrgRemoveDialogue
					orgName={org.friendly_name}
					onConfirm={() => {
						deleteOrg(org.id);
					}}
				>
					<Button
						onclick={(e: Event) => {
							e.preventDefault();
						}}
						size="icon"
						variant="ghost"
					>
						<Trash />
					</Button>
				</OrgRemoveDialogue>
			</CardAction>
		</CardHeader>

		<CardContent>
			<p class="text-sm">{org.description}</p>
			<div class="my-2 flex justify-between text-sm">
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
