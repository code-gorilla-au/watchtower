import type { PageLoad } from "./$types";
import { orgSvc, productSvc } from "$lib/watchtower";

export const load: PageLoad = async () => {
	await orgSvc.getDefault();
	const products = await productSvc.getAllByOrgId(orgSvc.defaultOrg?.id as number);

	return {
		organisation: orgSvc.defaultOrg,
		products
	};
};
