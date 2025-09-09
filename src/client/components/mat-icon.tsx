import { LitElement } from 'lit';
import { TemplateResult } from 'lit';

import { classMap } from 'lit/directives/class-map.js';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { html } from 'lit';
import { makeHostStyle } from '$client/components/icon-styles';
import { property } from 'lit/decorators.js';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'mat-icon': MatIcon;
  }
}

// 📘 display material icon

// 👉 https://marella.me/material-icons/demo/

//  --mat-icon-color    any color, default: inherit
//  --mat-icon-filter   any filter, default: none
//  --mat-icon-size     any size, default: 1em
//  --mat-icon-variant  outlined, round, sharp, two tone, default: (none)

@customElement('mat-icon')
export class MatIcon extends LitElement {
  static override styles = [
    makeHostStyle(),
    css`
      .material-icons {
        color: var(--mat-icon-color, inherit);
        direction: ltr;
        display: inline-block;
        filter: var(--mat-icon-filter, none);
        font-family: Material Icons;
        font-feature-settings: 'liga';
        font-size: var(--mat-icon-size, 1em);
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
        background-color: var(--mat-icon-color, inherit);
        color: transparent;
      }
    `
  ];

  @property() icon: string | null = null;

  override render(): TemplateResult {
    const style = getComputedStyle(this);
    const variant = style.getPropertyValue('--mat-icon-variant') ?? '';
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
