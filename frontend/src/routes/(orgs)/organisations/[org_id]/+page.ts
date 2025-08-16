import type { PageLoad } from "./$types";
import { orgSvc } from "$lib/watchtower";

export const load: PageLoad = async ({ params }) => {
	const org = await orgSvc.getById(Number(params.org_id));

	return { org };
};
