import type { PageLoad } from "./$types";
import { productSvc } from "$lib/watchtower";

export const load: PageLoad = async ({ params }) => {
	const product = await productSvc.getById(Number(params.product_id));
	const repos = await productSvc.getProductRepos(product.id);
	const prs = await productSvc.getOpenPrsByProduct(product.id);
	const securities = await productSvc.getSecurityByProduct(product.id);

	return { product, repos, prs, securities };
};
