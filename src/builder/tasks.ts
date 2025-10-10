import { bundle } from '$builder/bundler';
import { config } from '$builder/config';

// 📘 define all the tasks we can perform

class TaskClass {
  banner?: { color: string; icon: string } = {
    color: 'b7b8b4',
    icon: '\udb81\udcd8'
  };
  cmd?: string;
  cmds?: string[];
  description: string = '';
  func?: (args?: any) => Promise<boolean>;
  kill?: () => Promise<void>;
  name: string = '';
  subTasks?: string[];
  watchDirs?: string[];

  constructor(props: Task) {
    Object.assign(this, props);
  }
}

export interface Task extends TaskClass {}

const colors = {
  builder: '#639cf7',
  client: '#63f794',
  server: '#ebf763'
};

const icons = {
  assets: '',
  bundle: '',
  check: '',
  clean: '',
  css: '',
  js: '',
  lint: '',
  wasm: ''
};

// 👇 all the tasks we can perform

export const allTasks = [
  // ////////////////////////////////////////////////////////
  // 📘 BUNDLE
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'bundle:builder',
    description: 'Fully bundle builder',
    banner: { color: colors.builder, icon: icons.bundle },
    subTasks: ['check:builder', 'clean:builder', 'bundle:builder:js'],
    watchDirs: [
      config.paths.lib,
      `${config.paths.root}/tsconfig-app.json`,
      // 🔥 HACK -- a directory must come last
      config.paths['builder-ts']
    ]
  }),

  new TaskClass({
    name: 'bundle:builder:js',
    description: 'Bundle builder Javascript',
    banner: { color: colors.builder, icon: icons.js },
    func: ({ prod, verbose }): Promise<boolean> =>
      bundle({
        format: 'esm',
        outdir: `${config.paths['builder-js']}`,
        prod: !!prod,
        target: 'bun',
        verbose: !!verbose,
        roots: [`${config.paths['builder-ts']}/index.ts`],
        tsconfig: config.paths.tsconfig
      })
  }),

  new TaskClass({
    name: 'bundle:client',
    description: 'Fully bundle client',
    banner: { color: colors.client, icon: icons.bundle },
    subTasks: [
      'check:client',
      'clean:client',
      'test:client:wasm',
      'bundle:client:assets',
      'bundle:client:css',
      'bundle:client:js',
      'bundle:client:wasm'
    ],
    watchDirs: [
      config.paths.lib,
      `${config.paths.root}/tsconfig-app.json`,
      // 🔥 HACK -- a directory must come last
      config.paths['emulator-go'],
      config.paths['client-ts']
    ]
  }),

  new TaskClass({
    name: 'bundle:client:assets',
    description: 'Bundle client assets',
    banner: { color: colors.client, icon: icons.assets },
    cmds: [
      `cp ${config.paths['client-ts']}/index.html ${config.paths['client-js']}/`,
      `cp -r ${config.paths['client-ts']}/assets ${config.paths['client-js']}`
    ]
  }),

  new TaskClass({
    name: 'bundle:client:css',
    description: 'Bundle client CSS',
    banner: { color: colors.client, icon: icons.css },
    func: ({ prod, verbose }): Promise<boolean> =>
      bundle({
        outdir: `${config.paths['client-js']}`,
        prod: !!prod,
        target: 'bun',
        verbose: !!verbose,
        roots: [`${config.paths['client-ts']}/index.css`]
      })
  }),

  new TaskClass({
    name: 'bundle:client:js',
    description: 'Bundle client JS',
    banner: { color: colors.client, icon: icons.js },
    func: ({ prod, verbose }): Promise<boolean> =>
      bundle({
        format: 'esm',
        outdir: `${config.paths['client-js']}`,
        prod: !!prod,
        target: 'browser',
        verbose: !!verbose,
        roots: [`${config.paths['client-ts']}/index.ts`],
        tsconfig: config.paths.tsconfig
      })
  }),

  new TaskClass({
    name: 'bundle:client:wasm',
    description: 'Bundle client WASM',
    banner: { color: colors.client, icon: icons.wasm },
    cmd: `(cd ${config.paths['emulator-go']} && GOOS=js GOARCH=wasm go build -o ${config.paths['client-js']}/index.wasm main.go)`
  }),

  new TaskClass({
    name: 'bundle:server',
    description: 'Fully bundle server',
    banner: { color: colors.server, icon: icons.bundle },
    subTasks: ['check:server', 'clean:server', 'bundle:server:js'],
    watchDirs: [
      config.paths.lib,
      `${config.paths.root}/tsconfig-app.json`,
      // 🔥 HACK -- a directory must come last
      config.paths['server-ts']
    ]
  }),

  new TaskClass({
    name: 'bundle:server:js',
    description: 'Bundle server Javascript',
    banner: { color: colors.server, icon: icons.js },
    func: ({ prod, verbose }): Promise<boolean> =>
      bundle({
        format: 'esm',
        outdir: `${config.paths['server-js']}`,
        prod: !!prod,
        target: 'bun',
        verbose: !!verbose,
        roots: [`${config.paths['server-ts']}/index.ts`],
        tsconfig: config.paths.tsconfig
      })
  }),

  // ////////////////////////////////////////////////////////
  // 📘 CHECK
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'check',
    description: 'Check all code',
    subTasks: ['check:builder', 'check:client', 'check:server']
  }),

  new TaskClass({
    name: 'check:builder',
    description: 'Test compile builder without emitting JS',
    banner: { color: colors.builder, icon: icons.check },
    cmd: `bunx tsc --noEmit -p ${config.paths['builder-ts']}`
  }),

  new TaskClass({
    name: 'check:client',
    description: 'Test compile client without emitting JS',
    banner: { color: colors.client, icon: icons.check },
    cmd: `bunx tsc --noEmit -p ${config.paths['client-ts']}`
  }),

  new TaskClass({
    name: 'check:server',
    description: 'Test compile server without emitting JS',
    banner: { color: colors.server, icon: icons.check },
    cmd: `bunx tsc --noEmit -p ${config.paths['server-ts']}`
  }),

  // ////////////////////////////////////////////////////////
  // 📘 CLEAN
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'clean',
    description: 'Clean all code',
    subTasks: ['clean:builder', 'clean:client', 'clean:server']
  }),

  new TaskClass({
    name: 'clean:builder',
    description: 'Remove all files from builder dist',
    banner: { color: colors.builder, icon: icons.clean },
    cmds: [
      `mkdir -p ${config.paths['builder-js']}`,
      `touch ${config.paths['builder-js']}/at-least-one-to-rm`,
      `rm -rf ${config.paths['builder-js']}/*`
    ]
  }),

  new TaskClass({
    name: 'clean:client',
    description: 'Remove all files from client dist',
    banner: { color: colors.client, icon: icons.clean },
    cmds: [
      `mkdir -p ${config.paths['client-js']}`,
      `touch ${config.paths['client-js']}/at-least-one-to-rm`,
      `rm -rf ${config.paths['client-js']}/*`
    ]
  }),

  new TaskClass({
    name: 'clean:server',
    description: 'Remove all files from server dist',
    banner: { color: colors.server, icon: icons.clean },
    cmds: [
      `mkdir -p ${config.paths['server-js']}`,
      `touch ${config.paths['server-js']}/at-least-one-to-rm`,
      `rm -rf ${config.paths['server-js']}/*`
    ]
  }),

  // ////////////////////////////////////////////////////////
  // 📘 FORMAT
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'format',
    description: 'Format all code using prettier',
    cmd: `bunx prettier --write ${config.paths.root}`
  }),

  // ////////////////////////////////////////////////////////
  // 📘 LINT
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'lint',
    description: 'Lint all code using all available linters',
    subTasks: ['lint:eslint', 'lint:lit-analyzer', 'lint:stylelint']
  }),

  new TaskClass({
    name: 'lint:eslint',
    description: 'Lint build, client, and lib code with eslint',
    banner: { color: colors.builder, icon: icons.lint },
    cmd: `bunx eslint ${config.paths['builder-ts']} ${config.paths['client-ts']} ${config.paths['lib']} ${config.paths['server-ts']}`
  }),

  new TaskClass({
    name: 'lint:lit-analyzer',
    description: 'Lint client code using lit-analyzer',
    banner: { color: colors.client, icon: icons.lint },
    cmd: `bunx lit-analyzer ${config.paths['client-ts']}`
  }),

  new TaskClass({
    name: 'lint:stylelint',
    description:
      'Validate styles for CSS files and those embedded in TSX',
    banner: { color: colors.server, icon: icons.lint },
    cmd: `bunx stylelint --fix "${config.paths['client-ts']}/**/*.{css,tsx}"`
  }),

  // ////////////////////////////////////////////////////////
  // 📘 TEST
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'test',
    description: 'Run unit tests for all code',
    subTasks: ['test:client:wasm']
  }),

  new TaskClass({
    name: 'test:client:wasm',
    description: 'Run unit tests for client WASM',
    banner: { color: colors.client, icon: icons.wasm },
    // 🔥 can't test packages that depend on syscall/js
    cmd: `(cd ${config.paths['emulator-go']} && go test emulator/attrs emulator/buffer emulator/consts emulator/conv emulator/device emulator/glyph emulator/perf emulator/stack emulator/stream/inbound emulator/stream/outbound emulator/utils emulator/wcc -cover)`
  })
];

export const allTasksLookup: Record<string, Task> = allTasks.reduce(
  (acc, task) => {
    acc[task.name] = task;
    return acc;
  },
  <Record<string, Task>>{}
);
