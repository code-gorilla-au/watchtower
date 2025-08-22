<script lang="ts">
	import { Card, CardHeader, CardTitle } from "$components/ui/card/index.js";
	import { watchtower } from "$lib/wailsjs/go/models";
	import { Badge } from "$components/ui/badge/index.js";
	import { formatDate, truncate } from "$design/formats";

	type Props = {
		security: watchtower.SecurityDTO;
	};

	let { security }: Props = $props();

	function getSeverityVariant(
		severity: string
	): "default" | "secondary" | "destructive" | "outline" {
		switch (severity.toLowerCase()) {
			case "critical":
				return "destructive";
			case "high":
				return "destructive";
			case "medium":
				return "secondary";
			case "low":
				return "outline";
			default:
				return "default";
		}
	}

	function getStateVariant(state: string): "default" | "secondary" | "destructive" | "outline" {
		switch (state.toLowerCase()) {
			case "open":
				return "destructive";
			case "fixed":
				return "secondary";
			case "dismissed":
				return "outline";
			default:
				return "default";
		}
	}
</script>

<Card>
	<CardHeader class="flex items-center justify-between">
		<CardTitle>
			<span class="capitalize">{truncate(security.package_name)}</span>
		</CardTitle>
	</CardHeader>

	<div class="px-3">
		<div class="mb-2 flex items-center justify-between">
			<p class="text-sm text-muted-foreground">External ID</p>
			<p class="font-mono text-xs">{security.external_id}</p>
		</div>
		<div class="mb-2 flex items-center justify-between">
			<p class="text-sm text-muted-foreground">Severity</p>
			<Badge variant={getSeverityVariant(security.severity)}>{security.severity}</Badge>
		</div>
		<div class="mb-2 flex items-center justify-between">
			<p class="text-sm text-muted-foreground">State</p>
			<Badge variant={getStateVariant(security.state)}>{security.state}</Badge>
		</div>
		<div class="mb-2 flex items-center justify-between">
			<p class="text-sm text-muted-foreground">Patched Version</p>
			<p class="text-xs">{security.patched_version || "N/A"}</p>
		</div>
		<div class="mb-2 flex items-center justify-between">
			<p class="text-sm text-muted-foreground">Last updated</p>
			<p>{formatDate(security.updated_at)}</p>
		</div>
		<div>
			<Badge>{security.repository_name}</Badge>
		</div>
	</div>
</Card>
