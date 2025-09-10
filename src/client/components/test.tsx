import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { html } from 'lit';
import { property } from 'lit/decorators.js';
import { state } from 'lit/decorators.js';
import { stateContext } from '$client/state/state';

declare global {
  interface HTMLElementTagNameMap {
    'app-test': Test;
  }
}

// 📘 a test component

@customElement('app-test')
export class Test extends SignalWatcher(LitElement) {
  static override styles = css`
    :host {
      display: block;
    }
  `;

  @consume({ context: stateContext }) appState!: State;

  @state() job = 'dishwasher';

  @property() name = 'Bob';

  override render(): TemplateResult {
    return html`
      <p>As JSON ${this.appState.asJSON.get()}</p>
      <br />
      <br />
      <p>My name is ${this.name} and I am a ${this.job}</p>
    `;
  }
}
