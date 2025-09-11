import { Observable } from 'rxjs';
import { Observer } from 'rxjs';

import { e2a } from '$lib/convert';

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
          'color: palegreen',
          'color: skyblue'
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
          'color: palegreen',
          'color: coral'
        );
        observer.error(error);
      };
      // 👇 CLOSE
      this.#socket.onclose = (event: CloseEvent): void => {
        console.log(
          `%c3270 -> Server -> Client %cDisconnected ${event.reason}`,
          'color: palegreen',
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

  write(data: Uint8Array): void {
    this.#dump(data, 'Client -> Server -> 3270', true, 'yellow');
    this.#socket.send(data);
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
          'color: palegreen',
          'color: lightgray'
        );
        this.#socket.send(negotiator.encode(response));
      }
    } else {
      this.#dump(data, '3270 -> Server -> Client', true, 'palegreen');
      observer.next(data);
    }
  }

  #dump(
    data: Uint8Array,
    title: string,
    ebcdic = false,
    color = 'blue'
  ): void {
    const sliceSize = 32;
    let offset = 0;
    const total = data.length;
    console.groupCollapsed(
      `%c${title} ${ebcdic ? '(EBCDIC-encoded)' : ''}`,
      `color: ${color}`
    );
    console.log(
      '%c        00       04       08       0c       10       14       18       1c        00  04  08  0c  10  14  18  1c  ',
      'color: skyblue; font-weight: bold'
    );
    while (true) {
      const slice = new Uint8Array(
        data.slice(offset, Math.min(offset + sliceSize, total))
      );
      const { hex, str } = this.#dumpSlice(slice, sliceSize, ebcdic);
      console.log(
        `%c${this.#toHex(offset, 6)}: %c${hex} %c${str}`,
        'color: skyblue; font-weight: bold',
        'color: white',
        'color: wheat'
      );
      // setup for next time
      if (slice.length < sliceSize) break;
      offset += sliceSize;
    }
    console.groupEnd();
  }

  #dumpSlice(
    bytes: Uint8Array,
    sliceSize: number,
    ebcdic: boolean
  ): { hex: string; str: string } {
    let hex = '';
    let str = '';
    let ix = 0;
    // 👇 decode to hex and string equiv
    for (; ix < bytes.length; ix++) {
      const byte = bytes[ix];
      if (byte == null) break;
      hex += this.#toHex(byte, 2);
      const char = ebcdic ? e2a([byte]) : String.fromCharCode(byte);
      // NOTE: use special character in string as a visual aid to counting
      str += char === '\u00a0' || char === ' ' ? '\u2022' : char;
      if (ix > 0 && ix % 4 === 3) hex += ' ';
    }
    // 👇 pad remainder of slice
    for (; ix < sliceSize; ix++) {
      hex += '  ';
      str += ' ';
      if (ix > 0 && ix % 4 === 3) hex += ' ';
    }
    return { hex, str };
  }

  #toHex(num: number, pad: number): string {
    const padding = '0000000000000000'.substring(0, pad);
    const hex = num.toString(16);
    return padding.substring(0, padding.length - hex.length) + hex;
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
