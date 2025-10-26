import { Config } from '$client/state/state';
import { LitElement } from 'lit';
import { MdDialog } from '@material/web/dialog/dialog.js';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';
import { Tn3270 } from '$client/services/tn3270';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { globals } from '$client/css/globals/shadow-dom';
import { html } from 'lit';
import { query } from 'lit/decorators.js';
import { repeat } from 'lit/directives/repeat.js';
import { state } from 'lit/decorators.js';
import { stateContext } from '$client/state/state';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-connector': Connector;
  }
}

// 📘 manage 3270 connection

@customElement('app-connector')
export class Connector extends SignalWatcher(LitElement) {
  static override styles = [
    globals,
    css`
      .stretcher {
        align-items: center;
        display: flex;
        flex-direction: column;
        height: 100%;
        justify-content: center;
        min-width: 800px /* 👈 to prevent wrapping */;

        .connector {
          display: flex;
          flex-direction: column;
          gap: 2rem;
          justify-content: stretch;

          .header {
            text-align: center;

            .major {
              color: #00cc00;
              font-family: Terminal;
              font-size: 12rem;
              font-weight: bold;
              line-height: 1;
              text-shadow: 0.33rem 0.33rem #008000;
            }

            .minor {
              font-size: 2rem;
              font-weight: bold;
              text-transform: uppercase;
            }
          }

          .config {
            align-items: stretch;
            display: flex;
            flex-direction: row;
            gap: 1rem;
            justify-content: center;

            .connection,
            .emulation {
              display: flex;
              flex-direction: column;
              gap: 0.5rem;

              .instructions {
                font-weight: bold;
              }
            }

            .connection .controls {
              display: flex;
              flex-direction: row;
              gap: 0.5rem;
              justify-content: center;
            }

            .emulation .notes {
              font-size: smaller;
            }

            .label {
              align-items: center;
              display: flex;
              flex-direction: row;
              gap: 0.25rem;
            }
          }
        }
      }
    `
  ];

  @state() connecting!: boolean;
  @query('#dialog') dialog!: MdDialog;
  @query('#form') form!: HTMLFormElement;
  @state() message!: string;
  @consume({ context: stateContext }) state!: State;

  #dims: Record<string, [number, number]> = {
    '2': [24, 80],
    '3': [32, 80],
    '4': [43, 80],
    '5': [132, 27]
  };

  #tn3270: Tn3270 | null = null;

  async connect(evt: Event): Promise<void> {
    evt.preventDefault();
    const config = this.save();
    // 👇 now we can connect
    await this.connectImpl(config);
  }

  // 👁️ https://dev.to/blikblum/dry-form-handling-with-lit-19f
  // 👇 "connected" here means socket connection
  async connectImpl(config: Config): Promise<void> {
    try {
      this.connecting = true;
      this.#tn3270?.close();
      this.#tn3270 = await Tn3270.tn3270(
        config.host ?? 'localhost',
        config.port ?? '3270',
        `IBM-${config.device}-${config.model}-E`
      );
      this.#tn3270.stream$.subscribe({
        next: (bytes: Uint8ClampedArray) => {
          if (this.connecting)
            this.dispatchEvent(new CustomEvent('connected'));
          this.dispatchEvent(
            new CustomEvent('outbound', {
              detail: { bytes }
            })
          );
          this.connecting = false;
        },

        // 🔥 WebSocket connection established, but that failed
        error: async (e: any) => {
          console.error(e);
          this.connecting = false;
          this.message = e.reason;
          await this.dialog.show();
          this.dispatchEvent(new CustomEvent('disconnected'));
          this.#tn3270 = null;
        },

        // 👇 normal completion eg: Tn3270.close()
        complete: () => {
          console.log(
            `%c3270 -> Server -> Client %cClosed`,
            'color: palegreen',
            'color: cyan'
          );
          this.dispatchEvent(new CustomEvent('disconnected'));
          this.#tn3270 = null;
        }
      });
    } catch (e: any) {
      // 🔥 tried to upgrade to WebSocket, but that failed
      console.error(e);
      this.connecting = false;
      this.message = `Unable to reach proxy server ${location.hostname}:${location.port}`;
      await this.dialog.show();
      this.dispatchEvent(new CustomEvent('disconnected'));
      this.#tn3270 = null;
    }
  }

  debug(): void {
    // TODO 🔥 generalize to multiple test screenshots
    this.save('termtest');
    this.dispatchEvent(new CustomEvent('connected'));
  }

  disconnect(): void {
    this.#tn3270?.close();
  }

  panic(message: string): void {
    this.message = message;
    this.dialog.show();
    this.#tn3270?.close();
  }

  override render(): TemplateResult {
    return html`
      <main class="stretcher">
        <section class="connector">
          <article class="header">
            <header class="major">3270</header>
            <p class="minor">Go-powered 3270 Emulator</p>
          </article>

          <hr />

          <form
            @submit=${this.connect}
            class="config"
            id="form"
            name="config">
            <article class="connection">
              <div class="controls">
                <md-filled-text-field
                  label="Hostname or IP"
                  name="host"
                  style="width: 10rem"
                  value=${this.state.model.get().config
                    .host}></md-filled-text-field>
                <md-filled-text-field
                  label="Port"
                  name="port"
                  style="width: 5rem"
                  value=${this.state.model.get().config
                    .port}></md-filled-text-field>
              </div>

              <br />

              <div class="controls">
                <md-filled-button
                  ?disabled=${this.connecting}
                  style="align-self: center">
                  ${this.connecting ? 'Connnecting...' : 'Connect'}
                  <app-icon
                    icon="hourglass_full"
                    style=${styleMap({
                      display: this.connecting ? 'block' : 'none'
                    })}
                    slot="icon"></app-icon>
                </md-filled-button>

                <md-outlined-button @click=${this.setup} type="button">
                  Setup
                </md-outlined-button>

                <md-outlined-button @click=${this.debug} type="button">
                  Debug
                </md-outlined-button>
              </div>
            </article>

            <article class="emulation">
              <p class="instructions">Select 327x Model</p>

              ${repeat(
                Object.keys(this.#dims),
                (model) => model,
                (model) => html`
                  <label class="label">
                    <md-radio
                      ?checked=${this.state.model.get().config.model ===
                      model}
                      name="model"
                      value=${model}></md-radio>
                    Model ${model} &mdash;
                    <em class="notes">
                      ${this.#dims[model]?.[0]} x
                      ${this.#dims[model]?.[1]}
                    </em>
                  </label>
                `
              )}
            </article>

            <article class="emulation">
              <p class="instructions">Select 327x Device</p>

              <label class="label">
                <md-radio
                  ?checked=${this.state.model.get().config.device ===
                  '3278'}
                  name="device"
                  value="3278"></md-radio>
                3278 &mdash;
                <em class="notes">monochrome</em>
              </label>

              <label class="label">
                <md-radio
                  ?checked=${this.state.model.get().config.device ===
                  '3279'}
                  name="device"
                  value="3279"></md-radio>
                3279 &mdash;
                <em class="notes">color</em>
              </label>
            </article>
          </form>
        </section>
      </main>

      <md-dialog id="dialog">
        <header slot="headline">3270 Connection Error</header>
        <section slot="content">
          <p>
            An error occured while connecting to the terminal device at
            ${this.state.model.get().config
              .host}:${this.state.model.get().config.port}.
            Please take any necessary corrective action and retry.
            <br />
            <br />
            ${this.message}
          </p>
        </section>
        <form slot="actions" method="dialog">
          <md-outlined-button>OK</md-outlined-button>
        </form>
      </md-dialog>
    `;
  }

  save(screenshot = ''): Config {
    const formData = new FormData(this.form);
    const config = Object.fromEntries(
      formData.entries()
    ) as Partial<Config>;
    config.dims = this.#dims[config.model as string];
    config.screenshot = screenshot;
    this.state.updateConfig(config);
    return this.state.model.get().config;
  }

  sendToApp(bytes: Uint8ClampedArray): void {
    this.#tn3270?.sendToApp(bytes);
  }

  setup(): void {
    this.save();
    this.dispatchEvent(new CustomEvent('setup'));
  }
}
