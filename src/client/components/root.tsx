import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { Startup } from '$client/controllers/startup';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';

import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { html } from 'lit';
import { provide } from '@lit/context';
import { stateContext } from '$client/state/state';

declare global {
  interface HTMLElementTagNameMap {
    'app-root': Root;
  }
}

// 📘 the whole enchilada

@customElement('app-root')
export class Root extends SignalWatcher(LitElement) {
  static override styles = css`
    :host {
      display: block;
      margin: 1rem;
    }

    app-icon {
      /* --app-icon-color: palegreen; */
      /* --app-icon-filter: invert(8%) sepia(94%) saturate(4590%)
        hue-rotate(358deg) brightness(101%) contrast(112%); */
      --app-icon-size: 32px;
    }

    table {
      td {
        padding: 4px;
        text-align: left;
        vertical-align: middle;
      }
      td:first-child {
        text-align: center;
      }
    }
  `;

  @provide({ context: stateContext }) appState = new State('app-state');

  // eslint-disable-next-line no-unused-private-class-members
  #startup = new Startup(this);

  override render(): TemplateResult {
    const model = this.appState.model;
    return html`
      <section
        style="align-items: flex-start; display: flex; flex-direction: column; gap: 1rem">
        <label>
          <md-checkbox></md-checkbox>
          X is ${model.get().x}
        </label>

        <label>
          <md-checkbox></md-checkbox>
          Y is ${model.get().y}
        </label>

        <table>
          <tbody>
            <tr>
              <td>
                <app-icon icon="settings"></app-icon>
              </td>
              <td>Gear (material)</td>
            </tr>
          </tbody>
        </table>

        <md-filled-button
          @click=${(): void => this.appState.incrementX(10)}>
          Increment
        </md-filled-button>
      </section>

      <br />
      <br />
      <br />
      <br />
      <br />
      <app-test .name=${'Mark'}></app-test>
    `;
  }
}
