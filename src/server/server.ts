import { ServerWebSocket } from 'bun';
import { Socket as TCPSocket } from 'net';

import { cli } from '$server/cli';
import { log } from '$server/logger';
import { statSync } from 'node:fs';

let theServer: Bun.Server;

// 👇 we track context by a simple sessionID index

type Context = {
  host: string;
  port: string;
  sessionID: number;
  tcpSocket?: TCPSocket;
  webSocket?: ServerWebSocket<Context>;
};

const contexts = <Context[] | undefined[]>[];
let sessionID = 0;

// 👇 stop server and close all sockets on SIGINT

process.on('SIGINT', async () => {
  for (const ctx of contexts) {
    if (ctx) {
      log({
        important: `session ${ctx.sessionID}`,
        text: `... closing sockets`
      });
      ctx.tcpSocket?.end();
      ctx.webSocket?.close();
    }
  }
  log({
    important: `http://${theServer.hostname}:${theServer.port}`,
    text: `... stopping server`
  });
  await theServer.stop();
  process.exit(1);
});

// 👇 what does the CLI tell us to do?

const { port, root } = cli();

// 👇 server framework

async function serve(): Promise<void> {
  theServer = Bun.serve({
    fetch: fetchImpl,
    port,
    websocket: webSocketImpl
  });

  log({
    important: `http://${theServer.hostname}:${theServer.port}`,
    text: `starting server ...`
  });

  // 👇 never stop

  return new Promise(() => {});
}

// 👇 handle fetch requests for server - static content etc

const fetchImpl = (req: Request): any => {
  const url = new URL(req.url);
  if (url.pathname === '/') url.pathname = '/index.html';
  if (req.headers.get('upgrade') === 'websocket')
    return fetchUpgrade(url, req);
  else if (url.pathname === '/mtime') return fetchMTime();
  else return fetchStatic(url, req);
};

const fetchMTime = (): Response => {
  return new Response(String(statSync(root).mtimeMs));
};

const fetchStatic = (url: URL, req: Request): Response => {
  if (
    // 🔥 a quirk of Bun.serve ???
    !url.pathname.startsWith('/src:') &&
    !url.pathname.startsWith('/.')
  ) {
    log({ important: req.method, text: url.pathname });
    return new Response(Bun.file(`${root}${url.pathname}`));
  }
  // 👇 everything else is an error
  else {
    log({ warning: true, important: req.method, text: url.pathname });
    return new Response('Not found', { status: 404 });
  }
};

const fetchUpgrade = (url: URL, req: Request): Response | void => {
  if (
    theServer.upgrade(req, {
      data: {
        host: url.searchParams.get('host'),
        port: url.searchParams.get('port'),
        sessionID: sessionID++
      }
    })
  ) {
    return /* 👈 do not return a Response */;
  }
  return new Response('Failed to create socket connection to 3270', {
    status: 500
  });
};

// ↔️ handle socket SERVER / 3270

const tcpSocketImpl = (ctx: Context): void => {
  ctx.tcpSocket = new TCPSocket();
  ctx.tcpSocket.setNoDelay(true);
  // 👇 OPEN
  ctx.tcpSocket.connect(
    { host: ctx.host, port: parseInt(ctx.port) },
    () => {
      log({
        important: `#${ctx.sessionID} 3270 \uea99 SERVER`,
        text: `connected at ${ctx.host}:${ctx.port}`
      });
    }
  );
  // 👇 MESSAGE
  ctx.tcpSocket.on('data', (data: any) => {
    log({
      important: `#${ctx.sessionID} 3270 \uea99 SERVER`,
      text: 'forward message to CLIENT'
    });
    ctx.webSocket?.send(data);
  });
  // 👇 ERROR
  // 🔥 watch out - identical to close, except for log entry
  ctx.tcpSocket.on('error', (error: Error) => {
    log({
      error: true,
      important: `#${ctx.sessionID} 3270 \uea99 SERVER`,
      text: error.message
    });
    ctx.tcpSocket = undefined;
    ctx.webSocket?.close();
  });
  // 👇 CLOSE
  ctx.tcpSocket.on('end', () => {
    log({
      important: `#${ctx.sessionID} 3270 \uea99 SERVER`,
      text: 'disconnected'
    });
    ctx.tcpSocket = undefined;
    ctx.webSocket?.close();
  });
};

// ↔️ handle socket CLIENT / SERVER

const webSocketImpl = {
  // 👇 OPEN
  open: (ws: ServerWebSocket<Context>): void => {
    const ctx = { ...ws.data };
    contexts[ctx.sessionID] = ctx;
    ctx.webSocket = ws;
    tcpSocketImpl(ctx);
  },
  // 👇 MESSAGE
  message: (ws: ServerWebSocket<Context>, message: any): void => {
    const ctx = contexts[ws.data.sessionID];
    if (ctx) {
      log({
        important: `#${ctx.sessionID} CLIENT \uea99 SERVER`,
        text: 'forward message to 3270'
      });
      ctx.tcpSocket?.write(message);
    }
  },
  // 👇 ERROR
  // 🔥 watch out - identical to close, except for log entry
  error: (ws: ServerWebSocket<Context>, error: Error): void => {
    const ctx = contexts[ws.data.sessionID];
    if (ctx) {
      log({
        error: true,
        important: `#${ctx.sessionID} CLIENT \uea99 SERVER`,
        text: error.message
      });
      ctx.webSocket = undefined;
      ctx.tcpSocket?.end();
      contexts[ctx.sessionID] = undefined;
    }
  },
  // 👇 CLOSE
  close: (ws: ServerWebSocket<Context>): void => {
    const ctx = contexts[ws.data.sessionID];
    if (ctx) {
      log({
        important: `#${ctx.sessionID} CLIENT \uea99 SERVER`,
        text: 'disconnected'
      });
      ctx.webSocket = undefined;
      ctx.tcpSocket?.end();
      contexts[ctx.sessionID] = undefined;
    }
  }
};

// 👇 rock & roll

await serve();
