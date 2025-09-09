import { bundle } from '$builder/bundle';
import { config } from '$builder/config';
import { killServe } from '$builder/serve';
import { serve } from '$builder/serve';

// 📘 define all the tasks we can perform

class TaskClass {
  cmd?: string;
  cmds?: string[];
  description: string = '';
  func?: (args?: any) => Promise<any>;
  kill?: () => Promise<void>;
  name: string = '';
  subTasks?: string[];
  watchDirs?: string[];

  constructor(props: Task) {
    Object.assign(this, props);
  }
}

export interface Task extends TaskClass {}

// 👇 all the tasks we can perform

export const allTasks = [
  // ////////////////////////////////////////////////////////
  // 📘 BUNDLE
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'bundle:builder',
    description: 'Fully bundle builder',
    subTasks: ['clean:builder', 'check:builder', 'bundle:builder:js'],
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
    func: ({ prod, verbose }): Promise<any> =>
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
    cmds: [
      `mkdir -p ${config.paths['client-js']}`,
      `rm -rf ${config.paths['client-js']}/*`,
      `cp ${config.paths['client-ts']}/index.html ${config.paths['client-js']}/`,
      `cp ${config.paths['client-ts']}/favicon.ico ${config.paths['client-js']}/`,
      `cp -r ${config.paths['client-ts']}/assets ${config.paths['client-js']}`
    ],
    subTasks: ['check:client', 'bundle:client:css', 'bundle:client:js'],
    watchDirs: [
      config.paths.lib,
      `${config.paths.root}/tsconfig-app.json`,
      // 🔥 HACK -- a directory must come last
      config.paths['client-ts']
    ]
  }),
  new TaskClass({
    name: 'bundle:client:js',
    description: 'Bundle client Javascript',
    func: ({ prod, verbose }): Promise<any> =>
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
    name: 'bundle:client:css',
    description: 'Bundle client CSS',
    func: ({ prod, verbose }): Promise<any> =>
      bundle({
        outdir: `${config.paths['client-js']}`,
        prod: !!prod,
        target: 'bun',
        verbose: !!verbose,
        roots: [
          `${config.paths['client-ts']}/material-icons.css`,
          `${config.paths['client-ts']}/startup.css`
        ]
      })
  }),

  // ////////////////////////////////////////////////////////
  // 📘 CHECK
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'check',
    description: 'Check all code',
    subTasks: ['check:builder', 'check:client']
  }),

  new TaskClass({
    name: 'check:builder',
    description: 'Test compile builder without emitting JS',
    cmd: `bunx tsc --noEmit -p ${config.paths['builder-ts']}`
  }),

  new TaskClass({
    name: 'check:client',
    description: 'Test compile client without emitting JS',
    cmd: `bunx tsc --noEmit -p ${config.paths['client-ts']}`
  }),

  // ////////////////////////////////////////////////////////
  // 📘 CLEAN
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'clean',
    description: 'Clean all code',
    subTasks: ['clean:builder', 'clean:client']
  }),

  new TaskClass({
    name: 'clean:builder',
    description: 'Remove all files from builder dist',
    cmds: [
      `mkdir -p ${config.paths['builder-js']}`,
      `rm -rf ${config.paths['builder-js']}/*`
    ]
  }),

  new TaskClass({
    name: 'clean:client',
    description: 'Remove all files from client dist',
    cmds: [
      `mkdir -p ${config.paths['client-js']}`,
      `rm -rf ${config.paths['client-js']}/*`
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
    cmd: `bunx eslint ${config.paths['builder-ts']} ${config.paths['client-ts']} ${config.paths['lib']}`
  }),

  new TaskClass({
    name: 'lint:lit-analyzer',
    description: 'Lint client code using lit-analyzer',
    cmd: `bunx lit-analyzer ${config.paths['client-ts']}`
  }),

  new TaskClass({
    name: 'lint:stylelint',
    description:
      'Validate styles for CSS files and those embedded in TSX',
    cmd: `bunx stylelint --fix "${config.paths['client-ts']}/**/*.{css,tsx}"`
  }),

  // ////////////////////////////////////////////////////////
  // 📘 SERVE
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'serve',
    description: 'Launch test server',
    func: ({ verbose, watch }): Promise<any> =>
      serve({
        root: `${config.paths['client-js']}`,
        verbose: !!verbose,
        watch: !!watch
      }),
    kill: (): Promise<void> => killServe(),
    watchDirs: [config.paths['client-js']]
  })
];

export const allTasksLookup: Record<string, Task> = allTasks.reduce(
  (acc, task) => {
    acc[task.name] = task;
    return acc;
  },
  <Record<string, Task>>{}
);
