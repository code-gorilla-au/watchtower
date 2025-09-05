import { test, expect } from "@playwright/test";

test("route to register page", async ({ page }) => {
	await page.goto("/register/organisation");
	const title = page.getByText("Register - Organisation");
	await expect(title).toBeVisible();
});

test("should create organisation", async ({ page }) => {
	await page.goto("/register/organisation");
	await page.fill("input[id=friendly-name]", "Test Org");
	await page.fill("input[id=namespace]", "test-org");
	await page.fill("input[id=token]", "test-token");
	await page.fill("input[id=description]", "Test org description");
	await page.click("button[type=submit]");
	await expect(page).toHaveURL("/register/product");
});
