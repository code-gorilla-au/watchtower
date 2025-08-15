import { CreateOrganisation, GetDefaultOrganisation } from "$lib/wailsjs/go/watchtower/Service";
import { watchtower } from "$lib/wailsjs/go/models";
import { differenceInMinutes } from "date-fns";
import OrganisationDTO = watchtower.OrganisationDTO;
import { SvelteDate } from "svelte/reactivity";

export class OrgService {
	#internal: {
		defaultOrg?: OrganisationDTO;
		defaultLastSync?: Date;
		orgs?: OrganisationDTO[];
		orgsLastSync?: Date;
	};

	readonly defaultOrg: OrganisationDTO | undefined;

	constructor() {
		this.#internal = $state({
			defaultOrg: undefined,
			defaultLastSync: undefined,
			orgs: [],
			orgsLastSync: undefined
		});

		this.defaultOrg = $derived(this.#internal.defaultOrg);
	}

	async create(name: string, owner: string, token: string) {
		return await CreateOrganisation(name, owner, token);
	}

	async getDefault() {
		if (!this.defaultOrgStale()) {
			return this.defaultOrg;
		}

		const org = await GetDefaultOrganisation();
		this.#internal.defaultLastSync = new SvelteDate();
		this.#internal.defaultOrg = org;
		return this.defaultOrg;
	}

	private defaultOrgStale() {
		if (!this.#internal.defaultLastSync) {
			return true;
		}

		if (!this.#internal.defaultOrg) {
			return true;
		}

		const diff = differenceInMinutes(this.#internal.defaultLastSync, new SvelteDate());
		return diff > 5;
	}
}
