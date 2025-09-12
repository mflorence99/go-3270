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
  static override styles = [globals, css``];

  @consume({ context: stateContext }) theState!: State;

  override render(): TemplateResult {
    const model = this.theState.model;
    return html`
      <h1>Help #${model.get().pageNum}</h1>
    `;
  }
}
