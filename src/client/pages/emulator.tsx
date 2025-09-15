import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { globals } from '$client/css/globals/shadow-dom';
import { html } from 'lit';
import { stateContext } from '$client/state/state';

declare global {
  interface HTMLElementTagNameMap {
    'app-emulator': Emulator;
  }
}

// 📘 emulate the 3270 emulator

@customElement('app-emulator')
export class Emulator extends SignalWatcher(LitElement) {
  static override styles = [globals, css``];

  @consume({ context: stateContext }) theState!: State;

  override render(): TemplateResult {
    return html`
      <md-icon-button
        @click=${(): void => State.theTn3270?.close()}
        title="Disconnect from 3270">
        <app-icon icon="power_settings_new"></app-icon>
      </md-icon-button>
    `;
  }
}
