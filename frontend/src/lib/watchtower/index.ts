import { OrgService } from "$lib/watchtower/orgs.svelte";
import { ProductsService } from "$lib/watchtower/products.svelte";
import { NotificationsService } from "$lib/watchtower/notifications.svelte";
import { InsightsService } from "$lib/watchtower/insights.svelte";

const orgSvc = new OrgService();
const productSvc = new ProductsService();
const notificationSvc = new NotificationsService();
const insightsSvc = new InsightsService();

export { orgSvc, productSvc, notificationSvc, insightsSvc };
