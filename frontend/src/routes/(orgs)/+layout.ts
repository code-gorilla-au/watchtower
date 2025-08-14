import type { LayoutLoad } from "./$types";
import { OrgService, ProductsService } from "$lib/watchtower";
import { watchtower } from "$lib/wailsjs/go/models";
import ProductDTO = watchtower.ProductDTO;
import { goto } from "$app/navigation";

const org = new OrgService();
const productsSvc = new ProductsService();

export const load: LayoutLoad = async () => {
	let products: ProductDTO[] = [];
	const organisation = await org.getDefault();
	if (!organisation) {
		await goto("/register");
		return;
	}

	products = await productsSvc.getAllByOrgId(organisation?.id);

	return { organisation, products };
};
