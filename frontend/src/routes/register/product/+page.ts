import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { orgSvc } from "$lib/watchtower";

export const load: PageLoad = async () => {
	const organisation = await orgSvc.getDefault();
	if (!organisation) {
		error(404, "Organisation not found");
	}

	return { organisation };
};
