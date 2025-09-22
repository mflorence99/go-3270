import { Connector } from '$client/pages/connector';
import { Emulator } from '$client/pages/emulator';
import { LitElement } from 'lit';
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

// 🟧 shared by all pages

export type DataStreamEventDetail = {
  bytes: Uint8Array;
};

export const defaultColor = '#61b064';

export const Colors: Record<string, string> = {
  green: defaultColor,
  blue: '#42a5f5',
  orange: '#eb9a25',
  white: '#f9f9f9'
};

export const defaultDimensions: [number, number] = [80, 24];

export const Dimensions: Record<string, [number, number]> = {
  // 👇 [width, height]
  '1': [80, 12],
  '2': defaultDimensions,
  '3': [80, 32],
  '4': [80, 43],
  '5': [132, 27]
};

export const Emulators: Record<string, string> = {
  '1': 'IBM-3278-1',
  '2': 'IBM-3278-2',
  '3': 'IBM-3278-3',
  '4': 'IBM-3278-4',
  '5': 'IBM-3278-5'
};

export const Pages = {
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
  @query('.emulator') emulator!: Emulator;
  @state() pageNum = Pages.connector;
  @provide({ context: stateContext }) state = new State('state');

  // eslint-disable-next-line no-unused-private-class-members
  #startup = new Startup(this);

  override render(): TemplateResult {
    return html`
      <app-connector
        @connected=${(): any => (this.pageNum = Pages.emulator)}
        @datastream=${(e: CustomEvent<DataStreamEventDetail>): any =>
          this.emulator.datastream(e)}
        @disconnected=${(): any => (this.pageNum = Pages.connector)}
        class="connector"
        data-page-num="${Pages.connector}"
        style=${styleMap({
          opacity: this.pageNum === Pages.connector ? 1 : 0,
          zIndex: this.pageNum === Pages.connector ? 1 : -1
        })}></app-connector>

      <app-emulator
        @disconnect=${(): any => this.connector.disconnect()}
        @response=${(e: CustomEvent<DataStreamEventDetail>): any =>
          this.connector.response(e)}
        class="emulator"
        data-page-num="${Pages.emulator}"
        style=${styleMap({
          opacity: this.pageNum === Pages.emulator ? 1 : 0,
          zIndex: this.pageNum === Pages.emulator ? 1 : -1
        })}></app-emulator>
    `;
  }
}
