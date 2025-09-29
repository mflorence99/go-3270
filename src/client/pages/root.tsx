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
  #alarm = this.alarm.bind(this);
  #dumpBytes = this.dumpBytes.bind(this);
  #send = this.send.bind(this);

  // eslint-disable-next-line no-unused-private-class-members
  #startup = new Startup(this);

  async alarm(): Promise<void> {
    await this.dinger.play();
  }

  override connectedCallback(): void {
    super.connectedCallback();
    document.addEventListener('go3270-alarm', this.#alarm);
    document.addEventListener('go3270-dumpBytes', this.#dumpBytes);
    document.addEventListener('go3270-send', this.#send);
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback();
    document.removeEventListener('go3270-alarm', this.#alarm);
    document.removeEventListener('go3270-dumpBytes', this.#dumpBytes);
    document.removeEventListener('go3270-send', this.#send);
  }

  dumpBytes(evt: Event): void {
    const { bytes, title, ebcdic, color } = (evt as CustomEvent).detail;
    dumpBytes(bytes, title, ebcdic, color);
  }

  override render(): TemplateResult {
    return html`
      <app-connector
        @go3270-connected=${(): any => (this.pageNum = Pages.emulator)}
        @go3270-receive=${(evt: CustomEvent): any =>
          this.emulator.receive(evt.detail.bytes)}
        @go3270-disconnected=${(): any =>
          (this.pageNum = Pages.connector)}
        class="connector"
        data-page-num="${Pages.connector}"
        style=${styleMap({
          opacity: this.pageNum === Pages.connector ? 1 : 0,
          zIndex: this.pageNum === Pages.connector ? 1 : -1
        })}></app-connector>

      <app-emulator
        @disconnect=${(): any => this.connector.disconnect()}
        @go3270-send=${(evt: CustomEvent): any =>
          this.connector.send(evt.detail.bytes)}
        class="emulator"
        data-page-num="${Pages.emulator}"
        style=${styleMap({
          opacity: this.pageNum === Pages.emulator ? 1 : 0,
          zIndex: this.pageNum === Pages.emulator ? 1 : -1
        })}></app-emulator>

      <audio class="dinger" src="assets/ding.mp3"></audio>
    `;
  }

  send(evt: Event): void {
    const { bytes } = (evt as CustomEvent).detail;
    this.connector.send(bytes);
  }
}
