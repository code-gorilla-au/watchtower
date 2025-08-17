import type { PageLoad } from "./$types";
import { orgSvc } from "$lib/watchtower";

export const load: PageLoad = async () => {
	const org = await orgSvc.getDefault();

	return { org };
};
