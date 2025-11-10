import { config } from '$builder/config';
import { formatBytes } from '$lib/format';
import { log } from '$builder/logger';

type Params = {
  format?: 'iife' | 'esm' | 'cjs' | undefined;
  outdir: string;
  prod: boolean;
  roots: string[];
  target?: 'bun' | 'browser' | undefined;
  tsconfig?: string;
  verbose: boolean;
};

// 🟧 Run the Bub bundler

export async function bundle({
  format,
  outdir,
  prod,
  target,
  verbose,
  roots,
  tsconfig
}: Params): Promise<boolean> {
  // 👇 perform the build
  let build = null;
  try {
    build = await Bun.build({
      entrypoints: roots,
      format,
      minify: prod,
      outdir,
      sourcemap: true,
      target,
      tsconfig
    });
  } catch (e) {
    // 👇 TypeScript does not allow annotations on the catch clause
    const error = <AggregateError>e;
    // 👇 log any errors
    for (const message of error.errors) {
      log({
        error: message.error,
        warning: message.warning,
        text: `${message.name}-${message.message}`
      });
    }
  }

  // 👇 log any warnings
  if (verbose && build?.logs.length) {
    for (const message of build.logs) {
      log({
        warning: true,
        text: `${message.name}-${message.message}`
      });
    }
  }

  // 👇 log the artifacts
  if (verbose && build?.outputs.length) {
    for (const output of build.outputs) {
      const text = await output.text();
      log({
        important: formatBytes(text.length),
        text: config.makeRelative(output.path)
      });
    }
  }

  // 👇 we're done!
  return !!build?.success;
}
