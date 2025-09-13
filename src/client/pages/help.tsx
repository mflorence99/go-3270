import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { globals } from '$client/css/globals';
import { html } from 'lit';
import { stateContext } from '$client/state/state';

declare global {
  interface HTMLElementTagNameMap {
    'app-help': Help;
  }
}

// 📘 the whole enchilada

@customElement('app-help')
export class Help extends SignalWatcher(LitElement) {
  static override styles = [
    globals,
    css`
      section {
        display: grid;
        height: 100%;
        place-items: center;
        width: 100%;

        header {
          color: #00cc00;
          font-family: '3270 Font';
          font-size: 16rem;
          font-weight: bold;
          text-shadow: 0.33rem 0.33rem #008000;
        }
      }
    `
  ];

  @consume({ context: stateContext }) theState!: State;

  override render(): TemplateResult {
    // const model = this.theState.model;
    return html`
      <section>
        <header>3270</header>
      </section>
    `;
  }
}
