import type { PageLoad } from "./$types";
import { productSvc } from "$lib/watchtower";

export const load: PageLoad = async ({ parent }) => {
	const { organisation } = await parent();

	const products = await productSvc.getAllByOrgId(organisation?.id as number);

	return {
		products
	};
};
