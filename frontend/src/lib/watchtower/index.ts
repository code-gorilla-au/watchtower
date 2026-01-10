import { OrgService } from "$lib/watchtower/orgs.svelte";
import { ProductsService } from "$lib/watchtower/products.svelte";
import { NotificationsService } from "$lib/watchtower/notifications.svelte";
import { EventsOn } from "$lib/wailsjs/runtime";
import { EVENT_UNREAD_NOTIFICATIONS } from "$lib/watchtower/types";

const orgSvc = new OrgService();
const productSvc = new ProductsService();
const notificationSvc = new NotificationsService();

EventsOn(EVENT_UNREAD_NOTIFICATIONS, () => notificationSvc.getUnread());

export { orgSvc, productSvc, notificationSvc };
