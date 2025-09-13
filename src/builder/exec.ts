#!/usr/bin/env -S deno run --allow-all

import { Task } from '$builder/tasks';

import { $ } from 'bun';
import { allTasksLookup } from '$builder/tasks';
import { banner } from '$builder/logger';
import { cli } from '$builder/cli';
import { config } from '$builder/config';
import { debounce } from '$lib/debounce';
import { exit } from 'node:process';
import { kill } from 'node:process';
import { log } from '$builder/logger';

import chokidar from 'chokidar';
import psList from 'ps-list';

// 📘 execute all tasks to build & test

const { taskNames, prod, verbose, watch } = cli();

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

const todos: Task[] = reducer(taskNames ?? []);

// 👇 this closure will run each requested task

const error = (): void => {
  banner('See errors above', { color: '#ff8080', icon: '' });
  exit(1);
};

const run = async (todos: Task[]): Promise<void> => {
  for (const todo of todos) {
    try {
      // 👇 this looks pretty, but has no other function
      banner(todo.name, todo.banner);
      // 👇 could be a command
      const cmds = todo.cmds ?? [todo.cmd];
      for (const cmd of cmds) {
        if (cmd) {
          log({ important: todo.name, text: cmd });
          const plist = await psList();
          const existing = plist.find((p) => p.cmd === cmd);
          if (existing) kill(existing.pid, 'SIGINT');
          const { exitCode } = await $`${{ raw: cmd }}`.nothrow();
          if (exitCode !== 0) error();
        }
      }
      // 👇 could be a function
      if (todo.func) {
        log({ important: todo.name, text: 'function invoked' });
        await todo.kill?.();
        const result = await todo.func({ prod, verbose });
        if (!result) error();
      }
    } catch (e: any) {
      log({ error: true, data: e.message });
      error();
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
    const debounced = debounce(async () => {
      log({ warning: true, important: 'changes detected' });
      await run(todos);
    }, config.debounceMillis);
    // 👇 now create the watcher itself
    const watcher = chokidar.watch(allWatchedDirs, {
      persistent: true
    });
    log({ important: 'watching for changes', data: allWatchedDirs });
    watcher.on('all', (_, path) => {
      if (verbose) log({ important: 'changes detected', data: path });
      return debounced();
    });
    // 🔥 this hack trips the loop first time
    //    there'd better be a directory at the end of the list!!
    await $`touch ${allWatchedDirs.at(-1)}/.tickleme`;
  }
}

// 👇 otherwise, just run the todos and be done
if (allWatchedDirs.length === 0) {
  await run(todos);
  exit(0);
}
