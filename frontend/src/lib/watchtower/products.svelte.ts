import { watchtower } from "$lib/wailsjs/go/models";
import ProductDTO = watchtower.ProductDTO;
import { SvelteDate } from "svelte/reactivity";
import {
	CreateProduct,
	DeleteProduct,
	GetAllProductsForOrganisation,
	GetProductByID,
	GetProductPullRequests,
	GetProductRepos,
	GetPullRequestByOrganisation,
	GetSecurityByOrganisation,
	GetSecurityByProductID,
	SyncProduct,
	UpdateProduct
} from "$lib/wailsjs/go/watchtower/Service";
import RepositoryDTO = watchtower.RepositoryDTO;
import { differenceInMinutes } from "date-fns";
import { STALE_TIMEOUT_MINUTES } from "$lib/watchtower/types";

type RepoState = {
	data: RepositoryDTO[];
	lastSync?: SvelteDate;
};

export class ProductsService {
	#internal: {
		products: ProductDTO[];
		productsLastSync?: SvelteDate;
		repos: Record<number, RepoState>;
	};

	readonly products: ProductDTO[];

	constructor() {
		this.#internal = $state({
			products: [],
			repos: {}
		});

		this.products = $derived(this.#internal.products);
	}

	async create(name: string, description: string, orgId: number, tags: string[]) {
		const product = await CreateProduct(name, description, tags, orgId);
		this.#internal.products.push(product);

		return product;
	}

	async update(id: number, name: string, tags: string[]) {
		const product = await UpdateProduct(id, name, tags);
		this.internalUpdateProduct(product);

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
		if (!this.isProductStale()) {
			return this.internalGetProductById(id);
		}

		return this.getByIdForce(id);
	}

	async getProductRepos(productId: number) {
		if (!this.isRepoStale(productId)) {
			const repos = this.internalGetProductRepo(productId);
			if (!repos) {
				throw new Error("No repos for product found");
			}

			return repos.data;
		}

		const repos = await GetProductRepos(productId);
		this.#internal.repos[productId] = { data: repos, lastSync: new SvelteDate() };

		return repos;
	}

	async getOpenPrsByProduct(productId: number) {
		return await GetProductPullRequests(productId);
	}

	async getPullRequestsByOrganisation(orgId: number) {
		return await GetPullRequestByOrganisation(orgId);
	}

	async getSecurityByProduct(productId: number) {
		return await GetSecurityByProductID(productId);
	}

	async getSecurityByOrganisation(orgId: number) {
		return await GetSecurityByOrganisation(orgId);
	}

	async getAllByOrgId(orgId: number) {
		const products = await GetAllProductsForOrganisation(orgId);
		this.internalUpdateProducts(products);

		return products;
	}

	async syncProduct(id: number) {
		return await SyncProduct(id);
	}

	private async getByIdForce(id: number) {
		const product = await GetProductByID(id);
		this.internalUpdateProduct(product);
		return product;
	}

	private internalGetProductById(id: number) {
		const product = this.#internal.products.find((p) => p.id === id);
		if (!product) {
			throw new Error("Product not found");
		}

		return product;
	}

	private internalUpdateProduct(product: ProductDTO) {
		const idx = this.#internal.products.findIndex((p) => p.id === product.id);
		if (idx < 0) {
			return;
		}

		this.#internal.products.splice(idx, 1, product);
	}

	private internalUpdateProducts(products: ProductDTO[]) {
		this.#internal.products.splice(0, this.#internal.products.length, ...products);
		this.#internal.productsLastSync = new SvelteDate();
	}

	private isProductStale() {
		if (!this.#internal.productsLastSync) {
			return true;
		}

		if (!this.#internal.products.length) {
			return true;
		}

		const diff = differenceInMinutes(this.#internal.productsLastSync, new SvelteDate());
		return diff > STALE_TIMEOUT_MINUTES;
	}

	private internalGetProductRepo(productId: number) {
		if (!this.#internal.repos[productId]) {
			return undefined;
		}

		return this.#internal.repos[productId];
	}

	private isRepoStale(productId: number) {
		const repos = this.internalGetProductRepo(productId);
		if (!repos) {
			return true;
		}

		if (!repos.lastSync) {
			return true;
		}

		const diff = differenceInMinutes(repos.lastSync, new SvelteDate());
		return diff > STALE_TIMEOUT_MINUTES;
	}
}
