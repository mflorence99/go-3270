import { Signal } from '@lit-labs/signals';

import { computed } from '$client/types/signals';
import { config } from '$client/config';
import { createContext } from '@lit/context';
import { effect } from '$client/types/signals';
import { produce } from 'immer';
import { signal } from '@lit-labs/signals';

import StackTrace from 'stacktrace-js';

// 📘 base state class

abstract class Base<T> {
  // 👇 the signal that is the state itself
  model: Signal.State<T>;

  constructor(defaultState: T, key: string, persist: boolean) {
    if (persist) {
      const raw = localStorage.getItem(key);
      const persistedState = raw ? JSON.parse(raw) : defaultState;
      this.model = signal<T>(persistedState);
      effect(() =>
        localStorage.setItem(key, JSON.stringify(this.model.get()))
      );
    } else this.model = signal<T>(defaultState);
  }

  mutate(mutator: (state: T) => void): void {
    // 👇 finding the caller is "expensive" so we feature flag logging
    let caller: string | undefined;
    if (config.logStateChanges) {
      const frame = StackTrace.getSync()[1];
      caller = frame?.functionName;
    }
    // 👇 the "old" state
    const prevState = this.model.get();
    if (config.logStateChanges)
      console.log(
        '%c👈 prev state',
        'color: palegreen',
        caller,
        prevState
      );
    // 👇 the "new" state and (potentially) the patches that produced it
    const newState = produce(prevState, mutator, (patches) => {
      if (config.logStateChanges && patches)
        console.log(
          `%c🆕 patches... %c${caller} %c👉${JSON.stringify(patches)}`,
          'color: khaki',
          'color: white',
          'color: wheat'
        );
    });
    if (config.logStateChanges)
      console.log(
        '%c👉 next state',
        'color: skyblue',
        caller,
        newState
      );
    this.model.set(newState);
  }
}

// 📘 a conceptual model for real states
//    may morph into a real app-state

export type StateModel = {
  pageNum: number;
};

const defaultState: StateModel = {
  pageNum: 0
};

export class State extends Base<StateModel> {
  // 👇 just an example of a computed property
  asJSON = computed(() => JSON.stringify(this.model.get()));

  constructor(key: string) {
    super(defaultState, key, true);
  }

  // 👇 just an example of a mutator
  gotoPage(pageNum: number): void {
    this.mutate((state) => void (state.pageNum = pageNum));
  }

  // 👇 just an example of a mutator
  turnPage(): void {
    this.mutate((state) => {
      if (state.pageNum === 0) state.pageNum = 1;
      else if (state.pageNum === 1) state.pageNum = 2;
      else if (state.pageNum === 2) state.pageNum = 0;
    });
  }
}

export const stateContext = createContext<State>(Symbol('theState'));
