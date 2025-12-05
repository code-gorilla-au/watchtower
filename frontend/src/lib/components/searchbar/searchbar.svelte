<script lang="ts">
	import { truncate } from "$lib/hooks/formats";
	import { watchtower } from "$lib/wailsjs/go/models";
	import * as Command from "$lib/components/ui/command";
	import { OpenExternalURL } from "$lib/wailsjs/go/main/App";

	type Props = {
		prs: watchtower.PullRequestDTO[];
		securities: watchtower.SecurityDTO[];
		open: boolean;
	};

	let { prs, securities, open = $bindable() }: Props = $props();
</script>

<Command.Dialog bind:open>
	<Command.Input placeholder="Type a command or search..." />
	<Command.List>
		<Command.Empty>No results found.</Command.Empty>
		<Command.Group heading="Pull requests">
			{#each prs as pr (pr.id)}
				<Command.Item
					value={String(pr.id)}
					onSelect={() => {
						OpenExternalURL(pr.url);
					}}
					keywords={[pr.title, pr.repository_name, pr.author]}
				>
					<div class="flex flex-col">
						<span>{truncate(pr.title)}</span>
						<span class="text-sm text-muted-foreground"
							>{truncate(pr.repository_name)} - {pr.author}</span
						>
					</div>
				</Command.Item>
			{/each}
		</Command.Group>
		<Command.Separator />
		<Command.Group heading="Security vulnerability">
			{#each securities as s (s.id)}
				<Command.Item value={String(s.id)}>
					{s.package_name}
				</Command.Item>
			{/each}
		</Command.Group>
	</Command.List>
</Command.Dialog>
