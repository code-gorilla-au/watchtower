import { GetUnreadNotifications, MarkNotificationAsRead } from "$lib/wailsjs/go/watchtower/Service";

export class NotificationsService {
	/**
	 * Retrieves the unread notifications for the user.
	 */
	async getUnread() {
		return (await GetUnreadNotifications()) ?? [];
	}

	/**
	 * Marks a notification as read based on the provided notification ID.
	 */
	async markAsRead(id: number) {
		return await MarkNotificationAsRead(id);
	}

	/**
	 * Marks all unread notifications as read.
	 */
	async markAllAsRead() {
		const notifications = await this.getUnread();
		for (const notification of notifications) {
			await this.markAsRead(notification.id);
		}
	}
}
