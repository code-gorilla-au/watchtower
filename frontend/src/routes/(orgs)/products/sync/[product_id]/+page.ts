import type { PageLoad } from "./$types";
import { productSvc } from "$lib/watchtower";

export const load: PageLoad = async ({ params }) => {
	if (!params.product_id) {
	}
};
