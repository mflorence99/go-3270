import { ParseArgsOptionDescriptor } from 'node:util';

import { config } from '$server/config';
import { exit } from 'node:process';
import { log } from '$server/logger';
import { parseArgs } from 'node:util';
import { statSync } from 'node:fs';

// 📘 handle the details of the CLI

export type ParsedArgs = {
  help: boolean;
  port: string;
  root: string;
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
  port: {
    default: String(config.port),
    description: 'HTTP port',
    type: 'string'
  },
  root: {
    default: '',
    description: 'Root directory of client index.html',
    type: 'string'
  }
};

// 👇 launch the cli

export function cli(): ParsedArgs {
  let result;
  try {
    result = parseArgs({
      allowNegative: false,
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

  // 👇 log help data if requested

  if (parsedArgs.help) {
    logUsage();
    exit(0);
  }

  // 👇 validate the root path

  try {
    statSync(parsedArgs.root ?? '--root');
  } catch (e: any) {
    log({ error: true, important: e.message });
    exit(1);
  }

  return parsedArgs;
}

// 👇 log the CLI help data

function logUsage(msg = ''): void {
  console.log(`${msg}
Usage: server.ts [OPTION] ...

OPTIONS
-------
${logUsageOptions()}

    `);
}

function logUsageOptions(): string {
  return Object.entries(parseOptions)
    .reduce((acc, [name, options]) => {
      const k = `--${name}`.padEnd(10);
      const v = `${options.description}`;
      return `${acc}${k}${v}\n`;
    }, '')
    .trim();
}
