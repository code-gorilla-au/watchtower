export namespace insights {
	
	export class PullRequestInsights {
	    minDaysToMerge: number;
	    maxDaysToMerge: number;
	    avgDaysToMerge: number;
	    merged: number;
	    closed: number;
	    open: number;
	
	    static createFrom(source: any = {}) {
	        return new PullRequestInsights(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.minDaysToMerge = source["minDaysToMerge"];
	        this.maxDaysToMerge = source["maxDaysToMerge"];
	        this.avgDaysToMerge = source["avgDaysToMerge"];
	        this.merged = source["merged"];
	        this.closed = source["closed"];
	        this.open = source["open"];
	    }
	}
	export class SecurityInsights {
	    minDaysToFix: number;
	    maxDaysToFix: number;
	    avgDaysToFix: number;
	    fixed: number;
	    open: number;
	
	    static createFrom(source: any = {}) {
	        return new SecurityInsights(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.minDaysToFix = source["minDaysToFix"];
	        this.maxDaysToFix = source["maxDaysToFix"];
	        this.avgDaysToFix = source["avgDaysToFix"];
	        this.fixed = source["fixed"];
	        this.open = source["open"];
	    }
	}

}

export namespace notifications {
	
	export class Notification {
	    id: number;
	    organisation_id: number;
	    external_id: string;
	    status: string;
	    content: string;
	    type: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Notification(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.organisation_id = source["organisation_id"];
	        this.external_id = source["external_id"];
	        this.status = source["status"];
	        this.content = source["content"];
	        this.type = source["type"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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

export namespace organisations {
	
	export class OrganisationDTO {
	    id: number;
	    friendly_name: string;
	    description: string;
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
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.friendly_name = source["friendly_name"];
	        this.description = source["description"];
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
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class UpdateOrgParams {
	    ID: number;
	    FriendlyName: string;
	    Namespace: string;
	    Token: string;
	    Description: string;
	    DefaultOrg: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateOrgParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.FriendlyName = source["FriendlyName"];
	        this.Namespace = source["Namespace"];
	        this.Token = source["Token"];
	        this.Description = source["Description"];
	        this.DefaultOrg = source["DefaultOrg"];
	    }
	}

}

export namespace products {
	
	export class ProductDTO {
	    id: number;
	    name: string;
	    description: string;
	    tags: string[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ProductDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.tags = source["tags"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class PullRequestDTO {
	    id: number;
	    external_id: string;
	    title: string;
	    repository_name: string;
	    url: string;
	    state: string;
	    author: string;
	    tag: string;
	    product_name: string;
	    // Go type: time
	    merged_at: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PullRequestDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.external_id = source["external_id"];
	        this.title = source["title"];
	        this.repository_name = source["repository_name"];
	        this.url = source["url"];
	        this.state = source["state"];
	        this.author = source["author"];
	        this.tag = source["tag"];
	        this.product_name = source["product_name"];
	        this.merged_at = this.convertValues(source["merged_at"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class RepositoryDTO {
	    id: number;
	    name: string;
	    url: string;
	    topic: string;
	    owner: string;
	    product_name: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new RepositoryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.url = source["url"];
	        this.topic = source["topic"];
	        this.owner = source["owner"];
	        this.product_name = source["product_name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class SecurityDTO {
	    id: number;
	    external_id: string;
	    repository_name: string;
	    package_name: string;
	    state: string;
	    severity: string;
	    patched_version: string;
	    tag: string;
	    product_name: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SecurityDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.external_id = source["external_id"];
	        this.repository_name = source["repository_name"];
	        this.package_name = source["package_name"];
	        this.state = source["state"];
	        this.severity = source["severity"];
	        this.patched_version = source["patched_version"];
	        this.tag = source["tag"];
	        this.product_name = source["product_name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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

