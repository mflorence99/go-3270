import { Connector } from '$client/pages/connector';
import { Emulator } from '$client/pages/emulator';
import { LitElement } from 'lit';
import { MdDialog } from '@material/web/dialog/dialog.js';
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

// 🟧 The whole enchilada

@customElement('app-root')
export class Root extends SignalWatcher(LitElement) {
  static override styles = [
    globals,
    css`
      .connector,
      .emulator,
      .setup {
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

  @query('md-dialog[type="alert"]') alerter!: MdDialog;
  @query('.connector') connector!: Connector;
  @query('.emulator') emulator!: Emulator;
  @state() message = '';
  @state() pageNum = Pages.connector;
  @provide({ context: stateContext }) state = new State('state');

  mediator: Mediator;
  startup: Startup;

  constructor() {
    super();
    this.startup = new Startup(this);
    this.mediator = new Mediator(this);
  }

  alert(msg: string): void {
    this.message = msg;
    this.alerter.show();
  }

  flip(pageNum: number): void {
    this.pageNum = pageNum;
  }

  override render(): TemplateResult {
    return html`
      <app-connector
        @connected=${(): any => this.flip(Pages.emulator)}
        @disconnected=${(): any => this.flip(Pages.connector)}
        @outbound=${(evt: CustomEvent): any =>
          this.emulator.outbound(evt.detail.chars)}
        @setup=${(): any => (this.pageNum = Pages.setup)}
        class="connector"
        data-page-num="${Pages.connector}"
        style=${styleMap({
          opacity: this.pageNum === Pages.connector ? 1 : 0,
          zIndex: this.pageNum === Pages.connector ? 1 : -1
        })}></app-connector>

      <app-emulator
        @done=${(): any => (this.pageNum = Pages.connector)}
        class="emulator"
        data-page-num="${Pages.emulator}"
        style=${styleMap({
          opacity: this.pageNum === Pages.emulator ? 1 : 0,
          zIndex: this.pageNum === Pages.emulator ? 1 : -1
        })}></app-emulator>

      <app-setup
        @done=${(): any => this.flip(Pages.connector)}
        class="setup"
        data-page-num="${Pages.setup}"
        style=${styleMap({
          opacity: this.pageNum === Pages.setup ? 1 : 0,
          zIndex: this.pageNum === Pages.setup ? 1 : -1
        })}></app-setup>

      <md-dialog
        @closed=${(): void => this.mediator.alerterClosed()}
        @opened=${(): void => this.mediator.alerterOpened()}
        type="alert">
        <div slot="headline">Something Went Wrong</div>
        <form id="form" slot="content" method="dialog">
          ${this.message}
        </form>
        <div slot="actions">
          <md-filled-button form="form" value="ok">OK</md-filled-button>
        </div>
      </md-dialog>
    `;
  }
}
