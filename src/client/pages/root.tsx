import { Connector } from '$client/pages/connector';
import { Emulator } from '$client/pages/emulator';
import { LitElement } from 'lit';
import { Mediator } from '$client/controllers/mediator';
import { Pages } from '$client/controllers/mediator';
import { SignalWatcher } from '@lit-labs/signals';
import { Startup } from '$client/controllers/startup';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
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
  @query('.emulator') emulator!: Emulator;
  @state() pageNum = Pages.connector;
  @provide({ context: stateContext }) state = new State('state');

  // 👇 make sure "this" is right

  constructor() {
    super();
    new Startup(this);
    new Mediator(this);
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
    `;
  }
}
