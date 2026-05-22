import { test, expect } from "@playwright/test";

test("home renders DevForge", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "DevForge" })).toBeVisible();
});
