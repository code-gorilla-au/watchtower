import { watchtower } from "$lib/wailsjs/go/models";
import ProductDTO = watchtower.ProductDTO;
import { SvelteDate } from "svelte/reactivity";
import {
	CreateProduct,
	DeleteProduct,
	GetAllProductsForOrganisation,
	GetProductByID,
	GetProductRepos,
	SyncProduct,
	UpdateProduct
} from "$lib/wailsjs/go/watchtower/Service";
import RepositoryDTO = watchtower.RepositoryDTO;
import { differenceInMinutes } from "date-fns";

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

	async create(name: string, orgId: number, tags: string[]) {
		return CreateProduct(name, tags, orgId);
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
		this.#internal.products = products;
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
		return diff > 5;
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
		return diff > 5;
	}
}
