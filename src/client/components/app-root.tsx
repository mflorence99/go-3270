import { AppState } from '$client/state/app-state';
import { LitElement } from 'lit';
import { SignalWatcher } from '@lit-labs/signals';
import { StartupController } from '$client/controllers/startup';
import { TemplateResult } from 'lit';

import { appStateContext } from '$client/state/app-state';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { html } from 'lit';
import { provide } from '@lit/context';

declare global {
  interface HTMLElementTagNameMap {
    'app-root': AppRoot;
  }
}

// 📘 the whole enchilada

@customElement('app-root')
export class AppRoot extends SignalWatcher(LitElement) {
  static override styles = css`
    :host {
      display: block;
      margin: 1rem;
    }

    mat-icon {
      --mat-icon-color: palegreen;
      /* --mat-icon-filter: invert(8%) sepia(94%) saturate(4590%)
        hue-rotate(358deg) brightness(101%) contrast(112%); */
      --mat-icon-size: 32px;
      --mat-icon-variant: Two tone;
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

  @provide({ context: appStateContext }) appState = new AppState(
    'app-state'
  );

  // eslint-disable-next-line no-unused-private-class-members
  #startup = new StartupController(this);

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
                <mat-icon icon="settings"></mat-icon>
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
      <my-component .name=${'Mark'}></my-component>
    `;
  }
}
