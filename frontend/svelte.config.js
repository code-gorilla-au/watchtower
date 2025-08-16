import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter(),
		alias: {
			$design: "./src/lib/design",
			$components: "./src/lib/components"
		},
		prerender: {
			entries: [
				"*",
				"/register/organisation",
				"/register/product",
				"/products",
				"/organisations",
				"/organisations/[org_id]",
				"/products/create",
				"/products/[product_id]/sync",
				"/products/[product_id]/details"
			]
		}
	}
};

export default config;
