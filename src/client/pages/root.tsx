import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { Startup } from '$client/controllers/startup';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { globals } from '$client/css/globals';
import { html } from 'lit';
import { provide } from '@lit/context';
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
      app-help,
      app-home,
      app-screen {
        display: block;
        height: 100vh;
        opacity: 0;
        overflow: hidden;
        padding: 1rem;
        position: absolute;
        transition: opacity 0.5s ease-in-out;
        width: 100vw;
      }
    `
  ];

  @provide({ context: stateContext }) theState = new State('theState');

  // eslint-disable-next-line no-unused-private-class-members
  #startup = new Startup(this);

  override render(): TemplateResult {
    const model = this.theState.model;
    return html`
      <main @click=${(): void => this.theState.turnPage()}>
        <app-home
          style=${styleMap({
            opacity: model.get().pageNum === 0 ? 1 : 0
          })}></app-home>
        <app-help
          style=${styleMap({
            opacity: model.get().pageNum === 1 ? 1 : 0
          })}></app-help>
        <app-screen
          style=${styleMap({
            opacity: model.get().pageNum === 2 ? 1 : 0
          })}></app-screen>
      </main>
    `;
  }
}
