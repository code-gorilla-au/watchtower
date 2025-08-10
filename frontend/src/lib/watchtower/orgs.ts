import { CreateOrganisation, GetDefaultOrganisation } from "$lib/wailsjs/go/watchtower/Service";

export class OrgService {
	async create(name: string, owner: string) {
		return await CreateOrganisation(name, owner);
	}

	async getDefault() {
		return await GetDefaultOrganisation();
	}
}
