import { copyFile, mkdir, rm, stat } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const scriptPath = fileURLToPath(import.meta.url);
const defaultFrontendRoot = resolve(dirname(scriptPath), "..");
const assets = [
  ["swagger.html", "index.html"],
  ["swagger-ui-bundle.js", "swagger-ui-bundle.js"],
  ["swagger-ui-bundle.js.LICENSE.txt", "swagger-ui-bundle.js.LICENSE.txt"],
  ["swagger-ui.css", "swagger-ui.css"],
  ["LICENSE", "LICENSE"],
  ["NOTICE", "NOTICE"],
  ["favicon-16x16.png", "favicon-16x16.png"],
  ["favicon-32x32.png", "favicon-32x32.png"],
];

async function ensureFile(path) {
  try {
    await stat(path);
  } catch {
    throw new Error(`Required Swagger UI asset is missing: ${path}`);
  }
}

export async function buildApiDocs(frontendRoot = defaultFrontendRoot) {
  const sourceDirectory = resolve(frontendRoot, "documentation");
  const swaggerDirectory = resolve(
    frontendRoot,
    "node_modules/swagger-ui-dist",
  );
  const outputDirectory = resolve(frontendRoot, "public/docs");

  await rm(outputDirectory, { recursive: true, force: true });
  await mkdir(outputDirectory, { recursive: true });

  for (const [source, destination] of assets) {
    const base = source === "swagger.html" ? sourceDirectory : swaggerDirectory;
    const sourcePath = resolve(base, source);
    await ensureFile(sourcePath);
    await copyFile(sourcePath, resolve(outputDirectory, destination));
  }
}

if (process.argv[1] !== undefined && resolve(process.argv[1]) === scriptPath) {
  await buildApiDocs();
  console.log("Bundled local Swagger UI assets in public/docs.");
}
