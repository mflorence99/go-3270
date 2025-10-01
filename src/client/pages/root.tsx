import { Connector } from '$client/pages/connector';
import { Emulator } from '$client/pages/emulator';
import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { Startup } from '$client/controllers/startup';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { dumpBytes } from '$lib/dump';
import { globals } from '$client/css/globals/shadow-dom';
import { html } from 'lit';
import { provide } from '@lit/context';
import { query } from 'lit/decorators.js';
import { state } from 'lit/decorators.js';
import { stateContext } from '$client/state/state';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-root': Root;
  }
}

const Pages = {
  connector: 0,
  emulator: 1
};

// 📘 the whole enchilada

@customElement('app-root')
export class Root extends SignalWatcher(LitElement) {
  static override styles = [
    globals,
    css`
      .connector,
      .emulator {
        display: block;
        height: 100vh;
        opacity: 0;
        overflow: hidden;
        padding: 1rem;
        position: absolute;
        transition: opacity 0.5s ease-in-out;
        width: 100vw;
        z-index: -1;
      }
    `
  ];

  @query('.connector') connector!: Connector;
  @query('.dinger') dinger!: HTMLAudioElement;
  @query('.emulator') emulator!: Emulator;
  @state() pageNum = Pages.connector;
  @provide({ context: stateContext }) state = new State('state');

  // 👇 make sure "this" is right
  #disconnect = this.disconnect.bind(this);
  #go3270Message = this.go3270Message.bind(this);
  #keystroke = this.keystroke.bind(this);

  // eslint-disable-next-line no-unused-private-class-members
  #startup = new Startup(this);

  // 👇 "connected" here means DOM connection of this element
  override connectedCallback(): void {
    super.connectedCallback();
    // 👇 this comes from the Go code, requesting UI action
    document.addEventListener('go3270', this.#go3270Message);
    // 👇 these are pure UI events
    window.addEventListener('beforeunload', this.#disconnect);
    window.addEventListener('disconnect', this.#disconnect);
    window.addEventListener('keyup', this.#keystroke);
  }

  // 👇 "connected" here means socket connection to 3270
  disconnect(): void {
    this.connector.disconnect();
    this.emulator.disconnect();
  }

  // 👇 "connected" here means DOM connection of this element
  override disconnectedCallback(): void {
    super.disconnectedCallback();
    document.addEventListener('go3270', this.#go3270Message);
    window.removeEventListener('beforeunload', this.#disconnect);
    window.removeEventListener('disconnect', this.#disconnect);
    window.removeEventListener('keyup', this.#keystroke);
  }

  go3270Message(evt: Event): void {
    const params: Record<string, any> = (evt as CustomEvent).detail;
    switch (params.eventType) {
      case 'alarm':
        this.dinger.play();
        break;
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
          this.connector.sendToApp(bytes);
        }
        break;
    }
  }

  keystroke(evt: KeyboardEvent): void {
    if (this.pageNum === Pages.emulator) {
      const { altKey, code, ctrlKey, key, shiftKey } = evt;
      this.emulator.keystroke(code, key, altKey, ctrlKey, shiftKey);
      evt.preventDefault();
    }
  }

  override render(): TemplateResult {
    return html`
      <app-connector
        @connected=${(): any => (this.pageNum = Pages.emulator)}
        @disconnected=${(): any => (this.pageNum = Pages.connector)}
        @receiveFromApp=${(evt: CustomEvent): any =>
          this.emulator.receiveFromApp(evt.detail.bytes)}
        class="connector"
        data-page-num="${Pages.connector}"
        style=${styleMap({
          opacity: this.pageNum === Pages.connector ? 1 : 0,
          zIndex: this.pageNum === Pages.connector ? 1 : -1
        })}></app-connector>

      <app-emulator
        class="emulator"
        data-page-num="${Pages.emulator}"
        style=${styleMap({
          opacity: this.pageNum === Pages.emulator ? 1 : 0,
          zIndex: this.pageNum === Pages.emulator ? 1 : -1
        })}></app-emulator>

      <audio class="dinger" src="assets/ding.mp3"></audio>
    `;
  }
}
