import type { LayoutLoad } from "./$types";
import { OrgService } from "$lib/watchtower";
import { goto } from "$app/navigation";

const org = new OrgService();

export const load: LayoutLoad = async () => {
	const organisation = await org.getDefault();
	if (!organisation) {
		await goto("/register");
		return;
	}

	return { organisation };
};
