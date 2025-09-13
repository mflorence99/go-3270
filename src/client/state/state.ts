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
      // 👇 watch out! the new state maybe completely different from the old
      //    better to also remove any deprecated properties
      //    but this is simpler and always works, even for optionals
      this.model = signal<T>({ ...defaultState, ...persistedState });
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

// 📘 the entire state of the app

export type Config = {
  color: string;
  emulation: string;
  host: string;
  port: string;
};

export type StateModel = {
  config: Config;
  pageNum: number;
};

const defaultState: StateModel = {
  config: {
    color: 'green',
    emulation: '2',
    host: 'localhost',
    port: '3270'
  },
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

  // 👇 update the config
  updateConfig(config: Config): void {
    this.mutate((state) => void (state.config = config));
  }
}

export const stateContext = createContext<State>(Symbol('theState'));
