import { config } from '$builder/config';
import { log } from '$builder/logger';
import { statSync } from 'node:fs';

type Params = {
  root: string;
  verbose: boolean;
  watch: boolean;
};

// 📘 test server to deploy client

let server: Bun.Server;

export async function serve({
  root,
  verbose,
  watch
}: Params): Promise<any> {
  log({
    important: `http://localhost:${config.server.port}`,
    text: `starting test server ...`
  });

  // 👇 that's all we need!

  server = Bun.serve({
    fetch(req: Request): any {
      const url = new URL(req.url);
      let pathname = url.pathname;
      if (pathname === '/mtime')
        return new Response(String(statSync(root).mtimeMs));
      // 🔥 a quirk of Bun.serve ???
      else if (
        !pathname.startsWith('/src:') &&
        !pathname.startsWith('/.')
      ) {
        if (pathname === '/') pathname = '/index.html';
        if (verbose) log({ important: req.mode, text: pathname });
        return new Response(Bun.file(`${root}${pathname}`));
      } else return new Response('OK');
    },

    port: config.server.port
  });

  // 👇 stop server on SIGINT

  process.on('SIGINT', async () => {
    log({
      important: `http://localhost:${config.server.port}`,
      text: `... stopping test server`
    });
    await server.stop();
    process.exit();
  });

  // 👇 NOTE: watch mode does its own waiting

  return watch ? Promise.resolve() : new Promise(() => {});
}

export async function killServe(): Promise<void> {
  if (server) {
    log({
      important: `http://localhost:${config.server.port}`,
      text: `... stopping test server`
    });
    await server.stop();
  }
}
