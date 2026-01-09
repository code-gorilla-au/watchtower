import { describe, it, expect, vi, beforeEach } from "vitest";
import * as wails from "$lib/wailsjs/go/watchtower/Service";
import { notifications } from "$lib/wailsjs/go/models";
import { NotificationsService } from "$lib/watchtower/notifications.svelte";

describe("NotificationsService", () => {
	const spyGetUnread = vi.spyOn(wails, "GetUnreadNotifications");
	const spyMarkAsRead = vi.spyOn(wails, "MarkNotificationAsRead");

	beforeEach(() => {
		vi.clearAllMocks();
		vi.useFakeTimers();
	});

	describe("getUnread()", () => {
		beforeEach(() => {
			spyGetUnread.mockResolvedValue([
				new notifications.Notification({ id: 1, content: "Test notification" })
			]);
		});

		it("should return unread notifications", async () => {
			const notificationSvc = new NotificationsService();
			const list = await notificationSvc.getUnread();
			expect(list).toHaveLength(1);
			expect(spyGetUnread).toHaveBeenCalledTimes(1);
		});

		it("should have id and content", async () => {
			const notificationSvc = new NotificationsService();
			const list = await notificationSvc.getUnread();
			list.forEach((n) => {
				expect(n).toHaveProperty("id");
				expect(n).toHaveProperty("content");
			});
		});

		it("should not refetch if not stale", async () => {
			const notificationSvc = new NotificationsService();
			await notificationSvc.getUnread();
			expect(spyGetUnread).toHaveBeenCalledTimes(1);

			await notificationSvc.getUnread();
			expect(spyGetUnread).toHaveBeenCalledTimes(1);
		});

		it("should refetch if stale", async () => {
			const notificationSvc = new NotificationsService();
			await notificationSvc.getUnread();
			expect(spyGetUnread).toHaveBeenCalledTimes(1);

			vi.setSystemTime(new Date(Date.now() + 3 * 60 * 1000));

			await notificationSvc.getUnread();
			expect(spyGetUnread).toHaveBeenCalledTimes(2);
		});
	});

	describe("markAsRead()", () => {
		it("should call MarkNotificationAsRead and refresh list", async () => {
			spyMarkAsRead.mockResolvedValue();
			spyGetUnread.mockResolvedValue([]);

			const notificationSvc = new NotificationsService();
			await notificationSvc.markAsRead(1);

			expect(spyMarkAsRead).toHaveBeenCalledWith(1);
		});
	});

	describe("markAllAsRead()", () => {
		it("should mark each notification as read and refresh", async () => {
			const mockNotifications = [
				new notifications.Notification({ id: 1, content: "N1" }),
				new notifications.Notification({ id: 2, content: "N2" })
			];

			spyGetUnread.mockResolvedValueOnce(mockNotifications);
			spyMarkAsRead.mockResolvedValue();

			const notificationSvc = new NotificationsService();

			await notificationSvc.getUnread();

			spyGetUnread.mockResolvedValue([]);
			await notificationSvc.markAllAsRead();

			expect(spyMarkAsRead).toHaveBeenCalledTimes(2);
			expect(spyMarkAsRead).toHaveBeenCalledWith(1);
			expect(spyMarkAsRead).toHaveBeenCalledWith(2);
			expect(spyGetUnread).toHaveBeenCalledTimes(2);
		});
	});
});
