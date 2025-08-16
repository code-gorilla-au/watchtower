import type { LayoutLoad } from "./$types";
import { goto } from "$app/navigation";
import { orgSvc } from "$lib/watchtower";

export const load: LayoutLoad = async () => {
	await orgSvc.getDefault();
	if (!orgSvc.defaultOrg) {
		await goto("/register/organisation");
		return;
	}
};
