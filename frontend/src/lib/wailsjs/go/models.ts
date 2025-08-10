export namespace watchtower {
	export class OrganisationDTO {
		id: number;
		friendly_name: string;
		namespace: string;
		default_org: boolean;
		// Go type: time
		created_at: any;
		// Go type: time
		updated_at: any;

		static createFrom(source: any = {}) {
			return new OrganisationDTO(source);
		}

		constructor(source: any = {}) {
			if ("string" === typeof source) source = JSON.parse(source);
			this.id = source["id"];
			this.friendly_name = source["friendly_name"];
			this.namespace = source["namespace"];
			this.default_org = source["default_org"];
			this.created_at = this.convertValues(source["created_at"], null);
			this.updated_at = this.convertValues(source["updated_at"], null);
		}

		convertValues(a: any, classs: any, asMap: boolean = false): any {
			if (!a) {
				return a;
			}
			if (a.slice && a.map) {
				return (a as any[]).map((elem) => this.convertValues(elem, classs));
			} else if ("object" === typeof a) {
				if (asMap) {
					for (const key of Object.keys(a)) {
						a[key] = new classs(a[key]);
					}
					return a;
				}
				return new classs(a);
			}
			return a;
		}
	}
	export class ProductDTO {
		id: number;
		name: string;
		tags?: string;
		// Go type: time
		created_at: any;
		// Go type: time
		updated_at: any;

		static createFrom(source: any = {}) {
			return new ProductDTO(source);
		}

		constructor(source: any = {}) {
			if ("string" === typeof source) source = JSON.parse(source);
			this.id = source["id"];
			this.name = source["name"];
			this.tags = source["tags"];
			this.created_at = this.convertValues(source["created_at"], null);
			this.updated_at = this.convertValues(source["updated_at"], null);
		}

		convertValues(a: any, classs: any, asMap: boolean = false): any {
			if (!a) {
				return a;
			}
			if (a.slice && a.map) {
				return (a as any[]).map((elem) => this.convertValues(elem, classs));
			} else if ("object" === typeof a) {
				if (asMap) {
					for (const key of Object.keys(a)) {
						a[key] = new classs(a[key]);
					}
					return a;
				}
				return new classs(a);
			}
			return a;
		}
	}
}
