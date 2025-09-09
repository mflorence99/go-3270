import { cwd } from 'node:process';

// 📘 common configuration settings
//    NOT designed to be user-settable

const root = cwd();

export class ConfigClass {
  debounceMillis = 1000;

  paths = {
    'builder-js': `${root}/dist/builder`,
    'builder-ts': `${root}/src/builder`,
    'client-js': `${root}/dist/client`,
    'client-ts': `${root}/src/client`,
    'lib': `${root}/src/lib`,
    'root': root,
    'tsconfig': `${root}/tsconfig-app.json`
  };

  server = {
    port: 3000
  };

  relative(path: string): string {
    if (path.startsWith('/'))
      return path.substring(this.paths.root.length + 1);
    else return path;
  }
}

export const config: Readonly<ConfigClass> = new ConfigClass();
