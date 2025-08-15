import type { PageLoad } from "./$types";
import { productSvc } from "$lib/watchtower";

export const load: PageLoad = async ({ params }) => {
	const product = await productSvc.getById(Number(params.product_id));

	return { product };
};
