import { cli } from '$server/cli';
import { log } from '$server/logger';
import { statSync } from 'node:fs';

// 🔥 WIP

const { port, root } = cli();

async function serve(): Promise<void> {
  log({
    important: `http://localhost:${port}`,
    text: `starting server ...`
  });

  // 👇 that's all we need!

  const server = Bun.serve({
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
        log({ important: req.mode, text: pathname });
        return new Response(Bun.file(`${root}${pathname}`));
      } else return new Response('OK');
    },

    port: parseInt(port)
  });

  // 👇 stop server on SIGINT

  process.on('SIGINT', async () => {
    log({
      important: `http://localhost:${port}`,
      text: `... stopping server`
    });
    await server.stop();
    process.exit();
  });

  // 👇 never stop

  return new Promise(() => {});
}

await serve();
