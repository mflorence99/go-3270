import { parse } from 'node:path';

import chalk from 'chalk';
import figlet from 'figlet';
import StackTrace from 'stacktrace-js';

type Params = {
  data?: any;
  error?: boolean;
  important?: string;
  text?: string;
  warning?: boolean;
};

// 📘 provides a consistent logging format

export function log({
  data,
  error,
  important,
  text,
  warning
}: Params): void {
  // 👇 get the callers file, line # etc
  const frame = StackTrace.getSync()[1];
  const parsed = parse(frame?.fileName ?? 'unknown');
  // 👇 assemble the individual parts of the message
  const now = new Date().toLocaleTimeString();
  const parts: string[] = [
    error
      ? chalk.red.bold(now)
      : warning
        ? chalk.yellow.bold(now)
        : chalk.green(now),
    chalk.cyan(parsed.name).padEnd(13, '.')
  ];
  if (error) parts.push(chalk.red('🔥'));
  if (warning) parts.push(chalk.yellow('⚠️'));
  if (important) parts.push(chalk.yellow(important));
  if (text) parts.push(chalk.white(text));
  if (data) parts.push(chalk.cyan(JSON.stringify(data)));
  // 👇 ready to log them
  console.log(parts.join(' '));
}

// 📘 log short string using figlet

export function banner(str: string): void {
  console.log(
    chalk.green.bold(`\n\n  >>> ${str.toUpperCase().padEnd(72, ' ')}\n`)
  );
  console.log();
}

// 📘 log short string using figlet

export function figletize(str: string): void {
  console.log(
    chalk.green.bold(
      `${figlet.textSync(str.toUpperCase(), {
        font: 'Slant',
        horizontalLayout: 'fitted'
      })}`
    )
  );
}
