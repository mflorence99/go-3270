import { cwd } from 'node:process';

// 📘 common configuration settings
//    NOT designed to be user-settable

const root = cwd();

export class ConfigClass {
  debounceMillis = 250;

  keepAliveMillis = 250;

  paths = {
    'builder-js': `${root}/dist/builder`,
    'builder-ts': `${root}/src/builder`,
    'client-js': `${root}/dist/client`,
    'client-ts': `${root}/src/client`,
    'lib': `${root}/src/lib`,
    'root': root,
    'tsconfig': `${root}/tsconfig-app.json`,
    'server-js': `${root}/dist/server`,
    'server-ts': `${root}/src/server`
  };

  simulator = {
    http: {
      port: 8100
    },
    ws: {
      port: 8101
    }
  };
}

export const config: Readonly<ConfigClass> = new ConfigClass();
