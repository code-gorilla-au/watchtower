import {
	CreateOrganisation,
	DeleteOrganisation,
	GetAllOrganisations,
	GetDefaultOrganisation,
	GetOrganisationByID,
	SetDefaultOrg,
	UpdateOrganisation
} from "$lib/wailsjs/go/watchtower/Service";
import { watchtower } from "$lib/wailsjs/go/models";
import { differenceInMinutes } from "date-fns";
import OrganisationDTO = watchtower.OrganisationDTO;
import { SvelteDate } from "svelte/reactivity";

export class OrgService {
	#internal: {
		defaultOrg?: OrganisationDTO;
		defaultLastSync?: Date;
		orgs: OrganisationDTO[];
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
		const org = await CreateOrganisation(name, owner, token);
		this.internalUpdateOrg(org);
		if (org.default_org) {
			this.updateDefaultOrg(org);
		}

		return this.defaultOrg;
	}

	async update(params: { id: number; friendlyName: string; owner: string; defaultOrg: boolean }) {
		const updated = await UpdateOrganisation(
			new watchtower.UpdateOrgParams({
				ID: params.id,
				DefaultOrg: params.defaultOrg,
				FriendlyName: params.friendlyName,
				Namespace: params.owner
			})
		);

		this.internalUpdateOrg(updated);

		if (updated.default_org) {
			this.updateDefaultOrg(updated);
		}

		return updated;
	}

	async getById(id: number) {
		return GetOrganisationByID(id);
	}

	async getAll() {
		const orgs = await GetAllOrganisations();
		this.internalUpdateOrgs(orgs);

		return orgs;
	}

	async delete(id: number) {
		return DeleteOrganisation(id);
	}

	async getDefault() {
		if (!this.defaultOrgStale()) {
			return this.defaultOrg;
		}

		return this.getDefaultForce();
	}

	async setDefault(id: number) {
		const defaultOrg = await SetDefaultOrg(id);
		this.updateDefaultOrg(defaultOrg);
		return this.defaultOrg;
	}

	private async getDefaultForce() {
		const org = await GetDefaultOrganisation();
		this.updateDefaultOrg(org);

		return this.defaultOrg;
	}

	private updateDefaultOrg(org: OrganisationDTO) {
		this.#internal.defaultLastSync = new SvelteDate();
		this.#internal.defaultOrg = org;
	}

	private internalUpdateOrg(org: OrganisationDTO) {
		const idx = this.#internal.orgs?.findIndex((o) => o.id === org.id);
		if (idx < 0) {
			return;
		}

		this.#internal.orgs.splice(idx, 1, org);
	}

	private internalUpdateOrgs(orgs: OrganisationDTO[]) {
		this.#internal.orgs?.splice(0, this.#internal.orgs?.length, ...orgs);
		this.#internal.orgsLastSync = new SvelteDate();
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
