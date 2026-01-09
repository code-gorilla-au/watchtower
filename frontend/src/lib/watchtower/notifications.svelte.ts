import { notifications } from "$lib/wailsjs/go/models";
import { STALE_TIMEOUT_MINUTES } from "$lib/watchtower/types";
import { GetUnreadNotifications, MarkNotificationAsRead } from "$lib/wailsjs/go/watchtower/Service";

export class NotificationsService {
	readonly #notifications: notifications.Notification[];
	#lastSync?: number;

	constructor() {
		this.#notifications = $state([]);
	}
	/**
	 * Retrieves the unread notifications for the user.
	 */
	async getUnread() {
		if (this.isStale()) {
			await this.forceGetNotifications();
		}

		return this.#notifications;
	}

	/**
	 * Marks a notification as read based on the provided notification ID.
	 */
	async markAsRead(id: number) {
		await MarkNotificationAsRead(id);
		const idx = this.#notifications.findIndex((n) => n.id === id);
		if (idx < 0) {
			return;
		}

		this.#notifications.splice(
			idx,
			1,
			new notifications.Notification({
				...this.#notifications[idx],
				status: "read"
			})
		);
	}

	/**
	 * Marks all unread notifications as read.
	 */
	async markAllAsRead() {
		const notifications = [...(await this.getUnread())];
		for (const notification of notifications) {
			await this.markAsRead(notification.id);
		}

		await this.forceGetNotifications();
	}

	/**
	 * Forces the retrieval of unread notifications and updates the internal state with the fetched notifications.
	 */
	private async forceGetNotifications() {
		const notifications = (await GetUnreadNotifications()) ?? [];
		this.updateInternalNotifications(notifications);
	}

	/**
	 * Updates the internal notifications by replacing the current list with the provided notifications array.
	 */
	private updateInternalNotifications(notifications: notifications.Notification[]) {
		this.#notifications.splice(0, this.#notifications.length, ...notifications);
		this.#lastSync = Date.now();
	}

	/**
	 * Determines whether the current state is considered stale based on the last synchronisation time
	 * and the presence of notifications.
	 */
	private isStale() {
		if (!this.#lastSync || this.#notifications.length === 0) {
			return true;
		}

		const diff = (Date.now() - this.#lastSync) / (1000 * 60);
		return diff > STALE_TIMEOUT_MINUTES;
	}
}
