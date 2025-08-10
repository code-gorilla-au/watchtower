import type { PageLoad } from "./$types";
import { goto } from "$app/navigation";

export const load: PageLoad = async () => {
	await goto("/register");
};
