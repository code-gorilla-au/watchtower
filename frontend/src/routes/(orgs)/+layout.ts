import type { LayoutLoad } from "./$types";
import { goto } from "$app/navigation";
import { orgSvc } from "$lib/watchtower";
import { resolve } from "$app/paths";

export const load: LayoutLoad = async () => {
	await orgSvc.getAll();
	await orgSvc.getDefault();
	if (!orgSvc.defaultOrg) {
		await goto(resolve("/register"));
		return;
	}
};
