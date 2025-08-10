import { CreateOrganisation, GetDefaultOrganisation } from "$lib/wailsjs/go/watchtower/Service";
import { watchtower } from "$lib/wailsjs/go/models";
import OrganisationDTO = watchtower.OrganisationDTO;
import { differenceInMinutes } from "date-fns";

export class OrgService {
	#defaultOrg?: OrganisationDTO;
	#lastSync?: Date;

	async create(name: string, owner: string) {
		return await CreateOrganisation(name, owner);
	}

	async getDefault() {
		if (!this.#isStale() && this.#defaultOrg) {
			return this.#defaultOrg;
		}

		const org = await GetDefaultOrganisation();
		this.#lastSync = new Date();
		this.#defaultOrg = org;
		return org;
	}

	#isStale() {
		if (!this.#lastSync) {
			return true;
		}

		const diff = differenceInMinutes(this.#lastSync, new Date());
		return diff > 5;
	}
}
