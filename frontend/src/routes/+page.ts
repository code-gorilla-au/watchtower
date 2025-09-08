import type { PageLoad } from "./$types";
import { goto } from "$app/navigation";
import { orgSvc } from "$lib/watchtower";
import { resolve } from "$app/paths";

export const load: PageLoad = async () => {
	try {
		await orgSvc.getDefault();
		await goto(resolve("/dashboard"));
	} catch (e) {
		console.error(e);
		await goto(resolve("/register"));
	}
};
