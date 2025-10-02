import { ReactiveController } from 'lit';
import { Root } from '$client/pages/root';
import { Status } from '$client/state/state';

import { dumpBytes } from '$lib/dump';

export const Pages = {
  connector: 0,
  emulator: 1
};

// 📘 Mediate between Go code and the UI

export class Mediator implements ReactiveController {
  host: Root;

  #alarm!: HTMLAudioElement;

  // 👇 make sure "this" is right
  #disconnect = this.disconnect.bind(this);
  #go3270Message = this.go3270Message.bind(this);
  #keystroke = this.keystroke.bind(this);

  constructor(host: Root) {
    (this.host = host).addController(this);
  }

  disconnect(): void {
    this.host.connector.disconnect();
    this.host.emulator.disconnect();
  }

  async go3270Message(evt: Event): Promise<void> {
    switch ((evt as CustomEvent).detail.eventType) {
      case 'dumpBytes':
        {
          const { bytes, title, ebcdic, color } = (evt as CustomEvent)
            .detail;
          dumpBytes(bytes, title, ebcdic, color);
        }
        break;
      case 'log':
        {
          const { args } = (evt as CustomEvent).detail;
          console.log(...args.flat());
        }
        break;
      case 'sendToApp':
        {
          const { bytes } = (evt as CustomEvent).detail;
          this.host.connector.sendToApp(bytes);
        }
        break;
      case 'status':
        {
          const status: Partial<Status> = (evt as CustomEvent).detail;
          this.host.state.updateStatus(status);
          if (status.alarm) {
            await this.#alarm.play();
            this.host.state.updateStatus({ alarm: false });
          }
        }
        break;
    }
  }

  hostConnected(): void {
    // 👇 create an audio element to sound alarm
    this.#alarm = document.createElement('audio');
    this.#alarm.src = 'assets/ding.mp3';
    // 👇 this comes from the Go code, requesting UI action
    window.addEventListener('go3270', this.#go3270Message);
    // 👇 these are pure UI events
    window.addEventListener('beforeunload', this.#disconnect);
    window.addEventListener('disconnect', this.#disconnect);
    window.addEventListener('keyup', this.#keystroke);
  }

  hostDisconnected(): void {
    this.#alarm.remove();
    window.addEventListener('go3270', this.#go3270Message);
    window.removeEventListener('beforeunload', this.#disconnect);
    window.removeEventListener('disconnect', this.#disconnect);
    window.removeEventListener('keyup', this.#keystroke);
  }

  keystroke(evt: KeyboardEvent): void {
    if (this.host.pageNum === Pages.emulator) {
      const { altKey, code, ctrlKey, key, shiftKey } = evt;
      this.host.emulator.keystroke(
        code,
        key,
        altKey,
        ctrlKey,
        shiftKey
      );
      evt.preventDefault();
    }
  }
}
