import type { PageLoad } from "./$types";
import { goto } from "$app/navigation";
import { orgSvc } from "$lib/watchtower";

export const load: PageLoad = async () => {
	try {
		await orgSvc.getDefault();
		await goto("/dashboard");
	} catch (e) {
		console.error(e);
		await goto("/register");
	}
};
