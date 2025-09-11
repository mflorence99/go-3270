import { ServerWebSocket } from 'bun';
import { Socket as TCPSocket } from 'net';

import { cli } from '$server/cli';
import { log } from '$server/logger';
import { statSync } from 'node:fs';

import chalk from 'chalk';

// 🔥 WIP

const { port, root } = cli();

async function serve(): Promise<void> {
  log({
    important: `http://localhost:${port}`,
    text: `starting server ...`
  });

  let theTCPSocket: TCPSocket | null;
  let theWebSocket: ServerWebSocket<unknown> | null;

  const theServer = Bun.serve({
    //
    // 👇 simple fetch handler for static client code

    fetch(req: Request): any {
      const url = new URL(req.url);
      let pathname = url.pathname;
      // 👇 provide a WebSocket connection
      if (req.headers.get('upgrade') === 'websocket') {
        if (theServer.upgrade(req)) {
          return; // do not return a Response
        }
        return new Response('Upgrade failed', { status: 500 });
      }
      // 👇 support API for client hot reload
      else if (pathname === '/mtime')
        return new Response(String(statSync(root).mtimeMs));
      // 👇 deploy static content
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

    // 👇 the port, of course

    port: parseInt(port),

    // 👇 proxy client WebSocket

    websocket: {
      open(ws) {
        const host = 'localhost';
        const port = 3270;
        theTCPSocket = new TCPSocket();
        theTCPSocket.on('data', (data: any) => ws.send(data));
        theTCPSocket.on('error', (error: Error) => {
          console.log(
            chalk.green('3270 -> HOST'),
            chalk.red(error.message)
          );
        });
        theTCPSocket.on('end', () => {
          console.log(
            chalk.green('3270 -> HOST'),
            chalk.cyan('Disconnected')
          );
        });
        theTCPSocket.setNoDelay(true);
        theTCPSocket.connect({ host, port }, () => {
          console.log(
            chalk.green('3270 -> HOST'),
            chalk.blue(`Connected at ${host}:${port}`)
          );
        });
        theWebSocket = ws;
      },
      close() {
        theTCPSocket?.end();
        theTCPSocket = null;
        theWebSocket = null;
      },
      message(ws, message) {
        theTCPSocket?.write(message);
      }
    }
  });

  // 👇 stop server on SIGINT

  process.on('SIGINT', async () => {
    log({
      important: `http://localhost:${port}`,
      text: `... stopping server`
    });
    theTCPSocket?.end();
    theWebSocket?.close();
    await theServer.stop();
    process.exit(1);
  });

  // 👇 never stop

  return new Promise(() => {});
}

await serve();
