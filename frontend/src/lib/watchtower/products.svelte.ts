import { watchtower } from "$lib/wailsjs/go/models";
import ProductDTO = watchtower.ProductDTO;
import type { SvelteDate } from "svelte/reactivity";
import {
	CreateProduct,
	DeleteProduct,
	GetAllProductsForOrganisation,
	GetProductByID,
	SyncProduct,
	UpdateProduct
} from "$lib/wailsjs/go/watchtower/Service";

export class ProductsService {
	#internal: {
		products: ProductDTO[];
		lastSync?: SvelteDate;
	};

	readonly products: ProductDTO[];

	constructor() {
		this.#internal = $state({
			products: [],
			lastSync: undefined
		});

		this.products = $derived(this.#internal.products);
	}

	async create(name: string, orgId: number, tags: string[]) {
		return CreateProduct(name, tags, orgId);
	}

	async update(id: number, name: string, tags: string[]) {
		const product = await UpdateProduct(id, name, tags);
		const idx = this.#internal.products.findIndex((p) => p.id === id);
		if (idx < 0) {
			return;
		}

		this.#internal.products[idx] = product;
		return product;
	}

	async delete(id: number) {
		await DeleteProduct(id);
		const idx = this.#internal.products.findIndex((p) => p.id === id);
		if (idx < 0) {
			return;
		}

		this.#internal.products.splice(idx, 1);
	}

	async getById(id: number) {
		return GetProductByID(id);
	}

	async getAllByOrgId(orgId: number) {
		const products = await GetAllProductsForOrganisation(orgId);
		this.#internal.products.splice(0, this.#internal.products.length, ...products);
		return products;
	}

	async syncProduct(id: number) {
		return await SyncProduct(id);
	}
}
