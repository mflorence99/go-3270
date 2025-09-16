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
import { stateContext } from '$client/state/state';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-root': Root;
  }
}

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
      app-connector,
      app-emulator {
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

  @provide({ context: stateContext }) theState = new State('theState');

  // eslint-disable-next-line no-unused-private-class-members
  #startup = new Startup(this);

  override render(): TemplateResult {
    return html`
      <app-connector
        data-page-num="${Pages.connector}"
        style=${styleMap({
          opacity: this.#isPage(Pages.connector) ? 1 : 0,
          zIndex: this.#isPage(Pages.connector) ? 1 : -1
        })}></app-connector>

      <app-emulator
        data-page-num="${Pages.emulator}"
        style=${styleMap({
          opacity: this.#isPage(Pages.emulator) ? 1 : 0,
          zIndex: this.#isPage(Pages.emulator) ? 1 : -1
        })}></app-emulator>
    `;
  }

  #isPage(pageNum: number): boolean {
    const model = this.theState.model;
    return model.get().pageNum === pageNum;
  }
}
