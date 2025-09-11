import { Observable } from 'rxjs';
import { Observer } from 'rxjs';

const lookup: Record<string, number> = {
  BINARY: 0,
  DO: 253,
  DONT: 254,
  EOR: 25,
  IAC: 255,
  SB: 250,
  SE: 240,
  TERMINAL_TYPE: 24,
  WILL: 251,
  WONT: 252
};

const reverse: Record<number, string> = Object.fromEntries(
  Object.entries(lookup).map(([k, v]) => [v, k])
);

// 🟧 3270 Telnet factory

// 🟧 Raw Telnet to 3270
//    @see https://tools.ietf.org/html/rfc1576
//    @see https://tools.ietf.org/html/rfc1647
//    @see http://users.cs.cf.ac.uk/Dave.Marshall/Internet/node141.html

export class Tn3270 {
  stream$: Observable<Uint8Array>;

  #socket: WebSocket;

  private constructor(
    public host: string,
    public port: string,
    public model: string
  ) {
    this.#socket = new WebSocket(
      `ws://${location.hostname}:${location.port}?host=${host}&port=${port}`
    );
    this.stream$ = new Observable((observer: Observer<Uint8Array>) => {
      // 👇 OPEN
      this.#socket.onopen = (): void => {
        console.log(
          `%c3270 -> Server -> Client %cConnected ${host}:${port}`,
          'color: green',
          'color: blue'
        );
      };
      // 👇 MESSAGE
      this.#socket.onmessage = async (
        event: MessageEvent
      ): Promise<void> => {
        const data = new Uint8Array(await event.data.arrayBuffer());
        this.#dataHandler(data, observer);
      };
      // 🔥 ERROR
      this.#socket.onerror = (error): void => {
        console.log(
          `%c3270 -> Server -> Client %c${error.type}`,
          'color: green',
          'color: red'
        );
        observer.error(error);
      };
      // 👇 CLOSE
      this.#socket.onclose = (event: CloseEvent): void => {
        console.log(
          `%c3270 -> Server -> Client %cDisconnected ${event.reason}`,
          'color: green',
          'color: cyan'
        );
        observer.complete();
      };
      // 👇 return cleanup function called when complete
      return (): any => this.#socket.close();
    });
  }

  static async tn3270(
    host: string,
    port: string,
    model: string
  ): Promise<Tn3270 | null> {
    try {
      // 👇 initialize the WebSocket protocol
      await fetch(`http://${location.hostname}:${location.port}`, {
        headers: {
          upgrade: 'websocket'
        },
        mode: 'no-cors'
      });
      return new Tn3270(host, port, model);
    } catch (e: any) {
      console.error(e.message);
      return Promise.resolve(null);
    }
  }

  write(bytes: Uint8Array): void {
    this.#socket.send(bytes);
  }

  #dataHandler(data: Uint8Array, observer: Observer<Uint8Array>): void {
    if (data[0] === lookup.IAC) {
      const negotiator = new Negotiator(data);
      let response;
      if (negotiator.matches(['IAC', 'DO', 'TERMINAL_TYPE']))
        response = ['IAC', 'WILL', 'TERMINAL_TYPE'];
      if (negotiator.matches(['IAC', 'DO', 'EOR']))
        response = ['IAC', 'WILL', 'EOR', 'IAC', 'DO', 'EOR'];
      if (negotiator.matches(['IAC', 'DO', 'BINARY']))
        response = ['IAC', 'WILL', 'BINARY', 'IAC', 'DO', 'BINARY'];
      if (negotiator.matches(['IAC', 'SB', 'TERMINAL_TYPE']))
        response = [
          'IAC',
          'SB',
          'TERMINAL_TYPE',
          '0x00',
          this.model,
          'IAC',
          'SE'
        ];
      // 👇 send response
      if (response) {
        console.log(
          `%cClient -> Server -> 3270 %c${negotiator.decode()}`,
          'color: yellow',
          'color: white'
        );
        console.log(
          `%c3270 -> Server -> Client %c${response}`,
          'color: green',
          'color: gray'
        );
        this.write(negotiator.encode(response));
      }
    } else observer.next(data);
  }
}

// 🟧 Negotiate Telnet connection 3270 <-> Host

class Negotiator {
  constructor(private data: Uint8Array) {}

  decode(): string[] {
    const commands: string[] = [];
    for (let ix = 0; ix < this.data.length; ix++) {
      const byte = this.data[ix] ?? 0;
      let decoded = reverse[byte];
      // 👇 decode anything not in lookup as 0xXX
      if (typeof decoded === 'undefined')
        decoded = `0x${byte < 16 ? '0' : ''}${byte.toString(16)}`;
      commands.push(decoded);
    }
    return commands;
  }

  encode(commands: string[]): Uint8Array {
    const raw = commands.reduce(
      (acc, command) => {
        const encoded = lookup[command];
        // 👇 leave raw numbers as is
        if (typeof command === 'number') acc.push(command);
        // 👇 convert hex strings to decimal
        else if (command.startsWith('0x'))
          acc.push(parseInt(command.substring(3), 16));
        // 👇 anything not in lookup is a string, so decode bytes
        else if (typeof encoded === 'undefined') {
          for (let ix = 0; ix < command.length; ix++)
            acc.push(command.charCodeAt(ix));
        } else acc.push(encoded);
        return acc;
      },
      <number[]>[]
    );
    return new Uint8Array(raw);
  }

  matches(commands: string[]): boolean {
    return commands.every((command, ix) => {
      return lookup[command] === this.data[ix];
    });
  }
}
