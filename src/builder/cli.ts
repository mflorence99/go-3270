import { ParseArgsOptionDescriptor } from 'node:util';

import { allTasks } from '$builder/tasks';
import { allTasksLookup } from '$builder/tasks';
import { config } from '$builder/config';
import { exit } from 'node:process';
import { flattenObject } from '$builder/flatten';
import { log } from '$builder/logger';
import { parseArgs } from 'node:util';
import { statSync } from 'node:fs';

// 🟧 Handle the details of the CLI

export type ParsedArgs = {
  help: boolean;
  prod: boolean;
  taskNames: string[];
  verbose: boolean;
  watch: boolean;
};

interface ExtendedParseOptions extends ParseArgsOptionDescriptor {
  description: string;
}

const parseOptions: Record<string, ExtendedParseOptions> = {
  help: {
    default: false,
    description: 'Show this help information',
    short: 'h',
    type: 'boolean'
  },
  prod: {
    default: false,
    description: 'Build for production',
    short: 'p',
    type: 'boolean'
  },
  verbose: {
    default: false,
    description: 'Explain what is being done',
    short: 'v',
    type: 'boolean'
  },
  watch: {
    default: false,
    description: 'Run eligible tasks in watch mode',
    short: 'w',
    type: 'boolean'
  }
};

// 👇 launch the cli

export function cli(): ParsedArgs {
  let result;
  try {
    result = parseArgs({
      allowNegative: true,
      allowPositionals: true,
      args: Bun.argv,
      options: parseOptions,
      strict: true,
      tokens: false
    });
  } catch (e: any) {
    logUsage(e.message);
    exit(1);
  }
  const parsedArgs: ParsedArgs = { ...(<ParsedArgs>result.values) };
  parsedArgs.taskNames = result.positionals.slice(2) ?? [];

  // 👇 log help data if requested

  if (parsedArgs.help) {
    logUsage();
    exit(0);
  }

  // 👇 log no tasks error

  if (parsedArgs.taskNames.length === 0) {
    logNoop();
    exit(0);
  }

  // 👇 validate the requested tasks

  const allTaskNamesSet = new Set(Object.keys(allTasksLookup));
  const invalidTasks = parsedArgs.taskNames.filter(
    (taskName) => !allTaskNamesSet.has(taskName)
  );
  if (invalidTasks.length > 0) {
    logUsage(`Unknown task(s) '${invalidTasks}' requested`);
    exit(1);
  }

  // 👇 validate all the paths in the config

  const failures = [];
  for (const path of Object.values(config.paths)) {
    try {
      statSync(path);
      // 🔥 don't know why we need this
      // 🔥 specifying no-unused-vars doesn't work!
      // eslint-disable-next-line
    } catch (ignored) {
      failures.push(path);
    }
  }

  if (failures.length > 0) {
    failures.forEach((failure) =>
      log({ warning: true, important: failure, text: 'not found!' })
    );
  }

  // 👇 echo the config

  if (parsedArgs.verbose) {
    Object.entries(flattenObject(config)).forEach((entry) =>
      log({ important: `${entry.at(0)}`, text: `${entry.at(1)}` })
    );
    log({ data: parsedArgs, important: 'args' });
  }

  // 👇 confirm production mode

  if (parsedArgs.prod) {
    const result = confirm('Production mode requested. Are you sure?');
    if (!result) exit(1);
  }

  return parsedArgs;
}

// 👇 log abbreviated help if no tasks specified
//    patterned after Linux behavior for cp etc

function logNoop(): void {
  console.log(`
exec.ts: no tasks specified
try 'exec.ts --help' for more information.
    `);
}

// 👇 log the CLI help data

function logUsage(msg = ''): void {
  console.log(`${msg}
Usage: exec.ts [OPTION]... TASK...

OPTIONS
-------
${logUsageOptions()}

TASKS
-----
${logUsageTasks()}

    `);
}

// 🔥 we can do better than the magic "pad" numbers
//    but all this is already overkill!

function logUsageOptions(): string {
  return Object.entries(parseOptions)
    .reduce((acc, [name, options]) => {
      const k = `-${options.short}, --${name}`.padEnd(20);
      const v = `${options.description}`;
      return `${acc}${k}${v}\n`;
    }, '')
    .trim();
}

function logUsageTasks(): string {
  return allTasks
    .reduce((acc, task) => {
      const k = task.name.padEnd(20);
      const v = task.description.padEnd(50);
      const w =
        task.watchDirs && task.watchDirs.length > 0
          ? '(watchable)'
          : '';
      return `${acc}${k}${v}${w}\n`;
    }, '')
    .trim();
}
