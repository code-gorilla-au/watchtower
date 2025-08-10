import type { LayoutLoad } from "./$types";
import { OrgService } from "$lib/watchtower";

const org = new OrgService();

export const load: LayoutLoad = async () => {
	const organisation = await org.getDefault();

	return { organisation };
};
