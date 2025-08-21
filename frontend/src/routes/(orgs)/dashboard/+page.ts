import type { PageLoad } from "./$types";
import { orgSvc, productSvc } from "$lib/watchtower";

export const load: PageLoad = async () => {
	await orgSvc.getDefault();
	const products = await productSvc.getAllByOrgId(orgSvc.defaultOrg?.id as number);
	const prs = await productSvc.getPullRequestsByOrganisation(orgSvc.defaultOrg?.id as number);
	const securities = await productSvc.getSecurityByOrganisation(orgSvc.defaultOrg?.id as number);

	return {
		organisation: orgSvc.defaultOrg,
		products,
		prs,
		securities
	};
};
