import { CSSResult } from 'lit';

import { css } from 'lit';

// 📘 we must steal the styles etc from the global stylesheet
//    but do this once, efficiently, so we can use them
//    simply as regular styles as needed

export function makeStyleMaps(
  keyframeRegex: RegExp,
  ruleRegex: RegExp
): [
  keyframes: Map<string | null, string>,
  styles: Map<string | null, string>
] {
  const keyframes = new Map<string | null, string>();
  const styles = new Map<string | null, string>();
  for (const styleSheet of Array.from(document.styleSheets)) {
    for (const cssRule of Array.from(styleSheet.cssRules)) {
      // 👇 all the FontAwesome keyframes
      if (cssRule instanceof CSSKeyframesRule) {
        const keyframe = cssRule;
        const match = keyframeRegex.exec(keyframe.name);
        if (match?.[1]) keyframes.set(match[1], cssRule.cssText);
      }
      // 👇 all the FontAwesome styles
      if (cssRule instanceof CSSStyleRule) {
        const style = cssRule;
        const match = ruleRegex.exec(style.selectorText);
        if (match?.[1]) styles.set(match[1], cssRule.cssText);
      }
    }
  }
  return [keyframes, styles];
}

// 🆕 host style used by all icon components

export function makeHostStyle(): CSSResult {
  return css`
    :host {
      display: inline-block;
      text-align: center;
      vertical-align: middle;
    }
  `;
}
