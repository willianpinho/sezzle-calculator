import { screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

afterEach(() => {
  document.body.innerHTML = "";
  vi.resetModules();
});

describe("application bootstrap", () => {
  it("mounts the calculator into the application root", async () => {
    document.body.innerHTML = '<div id="root"></div>';

    await import("./main");

    expect(
      await screen.findByRole("heading", { name: "Calculator" }),
    ).toBeInTheDocument();
  });

  it("fails clearly when the application root is absent", async () => {
    await expect(import("./main")).rejects.toThrow(
      "Application root is missing",
    );
  });
});
