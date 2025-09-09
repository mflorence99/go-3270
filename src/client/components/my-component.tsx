import { AppState } from '$client/state/app-state';
import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { TemplateResult } from 'lit';

import { appStateContext } from '$client/state/app-state';
import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { html } from 'lit';
import { property } from 'lit/decorators.js';
import { state } from 'lit/decorators.js';

declare global {
  interface HTMLElementTagNameMap {
    'my-component': MyComponent;
  }
}

// 📘 a test component

@customElement('my-component')
export class MyComponent extends SignalWatcher(LitElement) {
  static override styles = css`
    :host {
      display: block;
    }
  `;

  @consume({ context: appStateContext }) appState!: AppState;

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
