#!/usr/bin/env -S deno run --allow-all

import { Task } from '$builder/tasks';

import { $ } from 'bun';
import { allTasksLookup } from '$builder/tasks';
import { cli } from '$builder/cli';
import { config } from '$builder/config';
import { debounce } from '$lib/debounce';
import { exit } from 'node:process';
import { figletize } from '$builder/logger';
import { kill } from 'node:process';
import { log } from '$builder/logger';

import chokidar from 'chokidar';
import psList from 'ps-list';

// 📘 execute all tasks to build & test Lintel
//    eg: exec.ts -p -w -v stylelint prettier

const { taskNames, prod, tedious, verbose, watch } = cli();

// 👇 flatten all the tasks and their subtasks into a sequense of todos

const reducer = (taskNames: (string | number)[]): Task[] => {
  return taskNames.reduce(
    (acc, taskName) => {
      const task = allTasksLookup[taskName];
      if (task) {
        acc.push(task);
        if (task.subTasks) acc.push(...reducer(task.subTasks));
      }
      return acc;
    },
    <Task[]>[]
  );
};

const todos: Task[] = reducer(taskNames);

// 👇 this closure will run each requested task

const run = async (todos: Task[]): Promise<void> => {
  for (const todo of todos) {
    try {
      // 👇 this looks pretty, but has no other function
      figletize(todo.name);
      // 👇 could be a command
      const cmds = todo.cmds ?? [todo.cmd];
      for (const cmd of cmds) {
        if (cmd) {
          log({ important: todo.name, text: cmd });
          const plist = await psList();
          const existing = plist.find((p) => p.cmd === cmd);
          if (existing) kill(existing.pid, 'SIGINT');
          const { exitCode } = await $`${{ raw: cmd }}`.nothrow();
          if (exitCode !== 0) break;
        }
      }
      // 👇 could be a function
      if (todo.func) {
        log({ important: todo.name, text: 'function invoked' });
        await todo.kill?.();
        const result = await todo.func({ prod, tedious, verbose });
        if (!result) break;
      }
    } catch (e: any) {
      log({ error: true, data: e.message });
      if (!watch) exit(1);
    }
  }
};

// 👇 if in watch mode, lookout for changes and run todos

let allWatchedDirs = [];

if (watch) {
  // 👇 consolidate the directories to watch
  allWatchedDirs = Array.from(
    todos.reduce((acc, todo) => {
      for (const dir of todo.watchDirs ?? []) acc.add(dir);
      return acc;
    }, new Set<string>())
  );

  // 👇 setup a watcher for the consolidated watchDirs
  //    we'll run all the todos because after debouncing we
  //    don't really know what changed

  if (allWatchedDirs.length > 0) {
    // 👇 create a debounced function that's invoked on changes
    const debounced = debounce(async (path) => {
      log({ important: 'changes detected', data: path });
      await run(todos);
    }, config.debounceMillis);
    // 👇 now create the watcher itself
    const watcher = chokidar.watch(allWatchedDirs, {
      persistent: true
    });
    log({ important: 'watching for changes', data: allWatchedDirs });
    watcher.on('all', (_, path) => debounced(path));
    // 🔥 this hack trips the loop first time
    //    there'd better be a directory at the end of the list!!
    await $`touch ${allWatchedDirs.at(-1)}/.tickleme`;
    // 👇 then run it on each change
    // for await (const event of watcher) {
    //   if (!['any', 'access'].includes(event.eventType))
    //     debounced(event);
    // }
  }
}

// 👇 otherwise, just run the todos
else await run(todos);

// 👇 that's all she wrote!

if (allWatchedDirs.length === 0) exit(0);
