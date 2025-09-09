import { Signal } from '@lit-labs/signals';

import { effect as subtleEffect } from 'signal-utils/subtle/microtask-effect';

// 📘 let's pretend we have first-order signal functions
//    that lit-signal doesn't yet provide

export function computed<T>(computation: () => T): Signal.State<T> {
  // @ts-ignore
  return new Signal.Computed<T>(computation);
}

export const effect = subtleEffect;
