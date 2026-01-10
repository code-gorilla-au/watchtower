import { OrgService } from "$lib/watchtower/orgs.svelte";
import { ProductsService } from "$lib/watchtower/products.svelte";
import { NotificationsService } from "$lib/watchtower/notifications.svelte";

const orgSvc = new OrgService();
const productSvc = new ProductsService();
const notificationSvc = new NotificationsService();

export { orgSvc, productSvc, notificationSvc };
