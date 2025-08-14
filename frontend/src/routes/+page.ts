import type { PageLoad } from "./$types";
import { goto } from "$app/navigation";
import { OrgService } from "$lib/watchtower";

const orgs = new OrgService();

export const load: PageLoad = async () => {
	try {
		await orgs.getDefault();
		await goto("/dashboard");
	} catch {
		await goto("/register/organisation");
	}
};
