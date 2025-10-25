import * as esbuild from "esbuild";
import fs from "node:fs";
import path from "node:path";
import { sassPlugin } from "esbuild-sass-plugin";

const OUT_DIR = "public/";

try {
  const files = fs.readdirSync(OUT_DIR);

  files.forEach((file) => {
    if (!file.endsWith(".html")) {
      const filePath = path.join(OUT_DIR, file);
      fs.rmSync(filePath);
      console.log(`Removing: ${filePath}`);
    }
  });
} catch (err) {
  console.error(`Error removing files: ${err.message}`);
}

let ctx = await esbuild.context({
  entryPoints: ["web/app.tsx"],
  bundle: true,
  minify: true,
  sourcemap: true,
  loader: { ".htm": "file" },
  outdir: `${OUT_DIR}`,
  plugins: [sassPlugin()],
});

let { host, port } = await ctx.serve({
  servedir: OUT_DIR,
});
await ctx.watch();

console.log(`[serve] listening at http://${host}:${port}`);
