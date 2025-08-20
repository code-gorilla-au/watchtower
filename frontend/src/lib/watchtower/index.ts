import { OrgService } from "$lib/watchtower/orgs.svelte";
import { ProductsService } from "$lib/watchtower/products.svelte";

const orgSvc = new OrgService();
const productSvc = new ProductsService();

export { orgSvc, productSvc };
