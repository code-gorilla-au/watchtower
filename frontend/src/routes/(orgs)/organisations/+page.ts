import type { PageLoad } from "./$types";
import { orgSvc } from "$lib/watchtower";

export const load: PageLoad = async () => {
	const orgs = await orgSvc.getAll();

	return { orgs };
};
