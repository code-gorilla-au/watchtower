import { notifications } from "$lib/wailsjs/go/models";
import { STALE_TIMEOUT_MINUTES } from "$lib/watchtower/types";
import { GetUnreadNotifications, MarkNotificationAsRead } from "$lib/wailsjs/go/watchtower/Service";

export class NotificationsService {
	readonly #unread: notifications.Notification[];
	#lastSync?: number;

	constructor() {
		this.#unread = $state([]);
	}
	/**
	 * Retrieves the unread notifications for the user.
	 */
	async getUnread(force: boolean = false) {
		if (this.isStale() || force) {
			await this.forceGetUnread();
		}

		return this.#unread;
	}

	/**
	 * Marks a notification as read based on the provided notification ID.
	 */
	async markAsRead(id: number) {
		await MarkNotificationAsRead(id);
		const idx = this.#unread.findIndex((n) => n.id === id);
		if (idx < 0) {
			return;
		}

		this.#unread.splice(idx, 1);
	}

	/**
	 * Marks all unread notifications as read.
	 */
	async markAllAsRead() {
		const notifications = [...(await this.getUnread())];
		for (const notification of notifications) {
			await this.markAsRead(notification.id);
		}

		await this.forceGetUnread();
	}

	/**
	 * Forces the retrieval of unread notifications and updates the internal state with the fetched notifications.
	 */
	private async forceGetUnread() {
		const notifications = (await GetUnreadNotifications()) ?? [];
		this.updateInternalUnread(notifications);
	}

	/**
	 * Updates the internal notifications by replacing the current list with the provided notifications array.
	 */
	private updateInternalUnread(notifications: notifications.Notification[]) {
		this.#unread.splice(0, this.#unread.length, ...notifications);
		this.#lastSync = Date.now();
	}

	/**
	 * Determines whether the current state is considered stale based on the last synchronisation time
	 * and the presence of notifications.
	 */
	private isStale() {
		if (!this.#lastSync || this.#unread.length === 0) {
			return true;
		}

		const diff = (Date.now() - this.#lastSync) / (1000 * 60);
		return diff > STALE_TIMEOUT_MINUTES;
	}
}
