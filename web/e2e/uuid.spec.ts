import { test, expect } from "@playwright/test";

test("uuid page generates uuids via API", async ({ page }) => {
  await page.goto("/uuid");
  await page.getByRole("button", { name: "Generate" }).click();
  const list = page.getByTestId("uuid-results");
  await expect(list.locator("li")).toHaveCount(3);
});
