import { expect, test } from "@playwright/test";

test.describe("homepage", () => {
  test("renders the Go template", async ({ page }) => {
    await page.goto("/");

    await expect(page.getByRole("heading", { name: "Homepage" })).toBeVisible();
  });

  test("mounts the preact island with props from Go", async ({ page }) => {
    await page.goto("/");

    const island = page.locator("[data-island='signup-form']");

    // The heading and the plan options exist only in web/pages.go - if they
    // render, data-props survived html/template's attribute escaping.
    await expect(island.getByRole("heading", { name: "Create your account" })).toBeVisible();
    await expect(island.locator("#plan option")).toHaveText(["Free", "Pro"]);
  });

  test("island state reacts to input", async ({ page }) => {
    await page.goto("/");

    const submit = page.locator("[data-island='signup-form'] button[type='submit']");
    const email = page.locator("#email");

    await expect(submit).toBeDisabled();

    await email.fill("not-an-email");
    await expect(page.getByText("Enter a valid email address.")).toBeVisible();
    await expect(submit).toBeDisabled();

    await email.fill("nawa@example.com");
    await expect(page.getByText("Enter a valid email address.")).toBeHidden();
    await expect(submit).toBeEnabled();
  });

  test("basecoat dialog opens and closes", async ({ page }) => {
    await page.goto("/");

    const dialog = page.locator("dialog.dialog");
    // basecoat 1.0 makes <dialog> a zero-size positioning anchor, so assert on
    // the <article> the user actually sees, not on the dialog element.
    const content = dialog.locator("article");

    await expect(content).toBeHidden();

    await page.getByRole("button", { name: "Confirm Changes" }).click();
    await expect(dialog).toHaveAttribute("open");
    await expect(content).toBeVisible();

    await dialog.getByRole("button", { name: "Cancel" }).click();
    await expect(content).toBeHidden();
  });

  test("no console errors on load", async ({ page }) => {
    const errors = [];
    page.on("console", (msg) => msg.type() === "error" && errors.push(msg.text()));
    page.on("pageerror", (err) => errors.push(err.message));

    await page.goto("/");
    await page.locator("[data-island='signup-form'] form").waitFor();

    expect(errors).toEqual([]);
  });
});
