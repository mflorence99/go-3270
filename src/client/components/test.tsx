import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { globals } from '$client/css/globals';
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
  static override styles = [
    globals,
    css`
      :host {
        display: block;
      }
    `
  ];

  @consume({ context: stateContext }) appState!: State;

  @state() job = 'dishwasher';

  @property() name = 'Bob';

  override render(): TemplateResult {
    return html`
      <p>As JSON ${this.appState.asJSON.get()}</p>
      <br />
      <a href="https://google.com" target="_blank">Google me!</a>
      <br />
      <p style="font-family: '3270 Font'; font-size: 1.5rem">
        My name is ${this.name} and I am a ${this.job}
      </p>
    `;
  }
}
