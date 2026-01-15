import type { PageLoad } from "./$types";
import { orgSvc, insightsSvc } from "$lib/watchtower";

export const load: PageLoad = async () => {
	await orgSvc.getDefault();
	const orgId = orgSvc.defaultOrg?.id;

	if (orgId) {
		await insightsSvc.refresh(orgId);
	}

	return {
		organisation: orgSvc.defaultOrg,
		insights: {
			pr: insightsSvc.prInsights,
			sec: insightsSvc.secInsights,
			window: insightsSvc.window
		}
	};
};
