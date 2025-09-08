import { config } from '$builder/config';

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
    cmd: `bunx tsc --noEmit -p ${config.paths['builder-ts']}`
  }),

  new TaskClass({
    name: 'check:client',
    description: 'Test compile client without emitting JS',
    cmd: `bunx tsc --noEmit -p ${config.paths['client-ts']}`
  }),

  new TaskClass({
    name: 'check:server',
    description: 'Test compile server without emitting JS',
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
    cmds: [
      `rm -rf ${config.paths['builder-js']}`,
      `mkdir -p ${config.paths['builder-js']}`
    ]
  }),

  new TaskClass({
    name: 'clean:client',
    description: 'Remove all files from client dist',
    cmds: [
      `rm -rf ${config.paths['client-js']}`,
      `mkdir -p ${config.paths['client-js']}`
    ]
  }),

  new TaskClass({
    name: 'clean:server',
    description: 'Remove all files from server dist',
    cmds: [
      `rm -rf ${config.paths['server-js']}`,
      `mkdir -p ${config.paths['server-js']}`
    ]
  }),

  // ////////////////////////////////////////////////////////
  // 📘 FORMAT
  // ////////////////////////////////////////////////////////

  new TaskClass({
    name: 'format',
    description: 'Format all code using prettier',
    cmd: `bunx prettier --write ${config.paths.root}`,
    watchDirs: [config.paths['client-ts']]
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
    description: 'Lint build, client, lib, and server code with eslint',
    cmd: `bunx eslint ${config.paths['builder-ts']} ${config.paths['client-ts']} ${config.paths['lib']} ${config.paths['server-ts']}`
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
  })
];

export const allTasksLookup: Record<string, Task> = allTasks.reduce(
  (acc, task) => {
    acc[task.name] = task;
    return acc;
  },
  <Record<string, Task>>{}
);
