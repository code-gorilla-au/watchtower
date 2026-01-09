import type { PageLoad } from "./$types";
import { notificationSvc } from "$lib/watchtower";

export const load: PageLoad = async () => {
	const notifications = await notificationSvc.getUnread();
	return { notifications };
};
