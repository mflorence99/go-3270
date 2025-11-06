import { ReactiveController } from 'lit';
import { Root } from '$client/pages/root';
import { Status } from '$client/state/state';

export const Pages = {
  connector: 0,
  emulator: 1,
  setup: 2
};

// 🟧 Mediate between Go code and the UI

export class Mediator implements ReactiveController {
  host: Root;

  #alarm!: HTMLAudioElement;

  // 👇 make sure "this" is right
  #blur = this.focussed(false);
  #disconnect = this.disconnect.bind(this);
  #focus = this.focussed(true);
  #go3270Message = this.go3270Message.bind(this);
  #keystroke = this.keystroke.bind(this);

  constructor(host: Root) {
    (this.host = host).addController(this);
  }

  disconnect(): void {
    this.host.connector.disconnect();
    this.host.emulator.disconnect();
  }

  focussed(focus: boolean): (evt: Event) => void {
    return () => {
      if (this.host.pageNum === Pages.emulator) {
        this.host.emulator.focussed(focus);
      }
    };
  }

  go3270Message(evt: Event): void {
    switch ((evt as CustomEvent).detail.eventType) {
      case 'panic':
        {
          const { args } = (evt as CustomEvent).detail;
          this.#alarm.play();
          // TODO 🔥 can we make this a modal dialog?
          window.alert(args);
          this.disconnect();
        }
        break;
      case 'inbound':
        {
          const { chars } = (evt as CustomEvent).detail;
          this.host.connector.sendToApp(chars);
        }
        break;
      case 'status':
        {
          // TODO 🔥 if we don't delay this, initial cursor not shown
          setTimeout(async () => {
            const status: Partial<Status> = (evt as CustomEvent).detail;
            this.host.state.updateStatus(status);
            if (status.alarm) {
              await this.#alarm.play();
              this.host.state.updateStatus({ alarm: false });
            }
          }, 0);
        }
        break;
    }
  }

  hostConnected(): void {
    // 👇 create an audio element to sound alarm
    this.#alarm = document.createElement('audio');
    this.#alarm.src = 'assets/beep.wav';
    // 👇 this comes from the Go code, requesting UI action
    window.addEventListener('go3270', this.#go3270Message);
    // 👇 these are pure UI events
    window.addEventListener('beforeunload', this.#disconnect);
    window.addEventListener('blur', this.#blur);
    window.addEventListener('disconnect', this.#disconnect);
    window.addEventListener('focus', this.#focus);
    window.addEventListener('keydown', this.#keystroke);
  }

  hostDisconnected(): void {
    this.#alarm.remove();
    window.addEventListener('go3270', this.#go3270Message);
    window.removeEventListener('beforeunload', this.#disconnect);
    window.removeEventListener('blur', this.#blur);
    window.removeEventListener('disconnect', this.#disconnect);
    window.removeEventListener('keydown', this.#keystroke);
    window.removeEventListener('focus', this.#focus);
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
