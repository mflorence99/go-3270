import { LitElement } from 'lit';
import { TemplateResult } from 'lit';

import { classMap } from 'lit/directives/class-map.js';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { html } from 'lit';
import { property } from 'lit/decorators.js';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-icon': Icon;
  }
}

// 📘 display material icon

// 👉 https://marella.me/material-icons/demo/

//  --app-icon-color    any color, default: inherit
//  --app-icon-filter   any filter, default: none
//  --app-icon-size     any size, default: 1em
//  --app-icon-variant  outlined, round, sharp, two tone, default: (none)

@customElement('app-icon')
export class Icon extends LitElement {
  static override styles = [
    css`
      :host {
        display: inline-block;
        text-align: center;
        vertical-align: middle;
      }

      .material-icons {
        color: var(--app-icon-color, inherit);
        direction: ltr;
        display: inline-block;
        filter: var(--app-icon-filter, none);
        font-family: Material Icons;
        font-feature-settings: 'liga';
        font-size: var(--app-icon-size, 1em);
        font-style: normal;
        font-weight: normal;
        letter-spacing: normal;
        line-height: 1;
        text-rendering: optimizeLegibility;
        text-transform: none;
        white-space: nowrap;
        word-wrap: normal;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
      }

      /* 👇 https://github.com/material-components/material-components-web-react/issues/730 */
      .material-icons-two-tone {
        background-clip: text;
        -webkit-background-clip: text;
        background-color: var(--app-icon-color, inherit);
        color: transparent;
      }
    `
  ];

  @property() icon: string | null = null;

  override render(): TemplateResult {
    const style = getComputedStyle(this);
    const variant = style.getPropertyValue('--app-icon-variant') ?? '';
    const fontFamily = `Material Icons ${variant}`.trim();
    const isTwotone = variant.toLowerCase() === 'two tone';
    return html`
      <i
        class=${classMap({
          'material-icons': true,
          'material-icons-two-tone': isTwotone
        })}
        style=${styleMap({ 'font-family': fontFamily })}>
        ${this.icon}
      </i>
    `;
  }
}
