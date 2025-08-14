<script lang="ts">
	import * as Menubar from "$lib/components/ui/menubar/index.js";
	import { Avatar } from "$components/ui/avatar/index.js";

	type Props = {};

	let {}: Props = $props();
</script>

<Menubar.Root>
	<Menubar.Menu>
		<Menubar.Trigger>
			<Avatar></Avatar>
		</Menubar.Trigger>
		<Menubar.Content>
			<Menubar.Item>
				New Tab
				<Menubar.Shortcut>⌘T</Menubar.Shortcut>
			</Menubar.Item>
			<Menubar.Item>New Window</Menubar.Item>
			<Menubar.Separator />
			<Menubar.Item>Share</Menubar.Item>
			<Menubar.Separator />
			<Menubar.Item>Print</Menubar.Item>
		</Menubar.Content>
	</Menubar.Menu>
</Menubar.Root>
