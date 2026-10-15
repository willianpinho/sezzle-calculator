// @vitest-environment node
import { mkdtemp, mkdir, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { afterEach, describe, expect, it, vi } from "vitest";
import { buildApiDocs } from "./build-api-docs.mjs";

const assetNames = [
  "swagger-ui-bundle.js",
  "swagger-ui-bundle.js.LICENSE.txt",
  "swagger-ui.css",
  "LICENSE",
  "NOTICE",
  "favicon-16x16.png",
  "favicon-32x32.png",
];
const temporaryDirectories = [];
const scriptPath = fileURLToPath(
  new URL("./build-api-docs.mjs", import.meta.url),
);

async function createFixture() {
  const root = await mkdtemp(join(tmpdir(), "sezzle-api-docs-"));
  temporaryDirectories.push(root);
  await mkdir(join(root, "documentation"), { recursive: true });
  await mkdir(join(root, "node_modules/swagger-ui-dist"), { recursive: true });
  await writeFile(
    join(root, "documentation/swagger.html"),
    "<main>docs</main>",
  );
  await Promise.all(
    assetNames.map((name) =>
      writeFile(
        join(root, "node_modules/swagger-ui-dist", name),
        `asset:${name}`,
      ),
    ),
  );
  return root;
}

afterEach(async () => {
  await Promise.all(
    temporaryDirectories
      .splice(0)
      .map((path) => rm(path, { recursive: true, force: true })),
  );
});

describe("buildApiDocs", () => {
  it("does not run the CLI entry when imported as a helper", async () => {
    const root = await createFixture();
    const log = vi.spyOn(console, "log").mockImplementation(() => {});
    const module = await import("./build-api-docs.mjs?helper");

    expect(log).not.toHaveBeenCalled();
    await module.buildApiDocs(root);
    log.mockRestore();

    await expect(
      readFile(join(root, "public/docs/index.html"), "utf8"),
    ).resolves.toBe("<main>docs</main>");
  });

  it("runs the CLI entry only when Node executes this script", async () => {
    const originalArgument = process.argv[1];
    const log = vi.spyOn(console, "log").mockImplementation(() => {});
    process.argv[1] = resolve(scriptPath);

    try {
      await import("./build-api-docs.mjs?cli");
      expect(log).toHaveBeenCalledWith(
        "Bundled local Swagger UI assets in public/docs.",
      );
      await expect(
        readFile(
          join(dirname(scriptPath), "../public/docs/index.html"),
          "utf8",
        ),
      ).resolves.toContain("swagger-ui");
    } finally {
      process.argv[1] = originalArgument;
      log.mockRestore();
    }
  });

  it("does not throw or run when Node has no script argument", async () => {
    const originalArgument = process.argv[1];
    const log = vi.spyOn(console, "log").mockImplementation(() => {});
    process.argv[1] = undefined;

    try {
      await import("./build-api-docs.mjs?no-entry");
      expect(log).not.toHaveBeenCalled();
    } finally {
      process.argv[1] = originalArgument;
      log.mockRestore();
    }
  });

  it("copies local template, Swagger assets, and replaces stale output", async () => {
    const root = await createFixture();
    await mkdir(join(root, "public/docs"), { recursive: true });
    await writeFile(join(root, "public/docs/stale.txt"), "stale");

    await buildApiDocs(root);

    await expect(
      readFile(join(root, "public/docs/index.html"), "utf8"),
    ).resolves.toBe("<main>docs</main>");
    await expect(
      readFile(join(root, "public/docs/swagger-ui.css"), "utf8"),
    ).resolves.toBe("asset:swagger-ui.css");
    await expect(
      readFile(join(root, "public/docs/stale.txt"), "utf8"),
    ).rejects.toMatchObject({ code: "ENOENT" });
  });

  it("fails with the missing local asset path", async () => {
    const root = await createFixture();
    await rm(join(root, "node_modules/swagger-ui-dist/NOTICE"));

    await expect(buildApiDocs(root)).rejects.toThrow("swagger-ui-dist/NOTICE");
  });
});
