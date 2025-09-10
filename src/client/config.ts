// 📘 common configuration settings
//    NOT designed to be user-settable

export class ConfigClass {
  logStateChanges = location.hostname === 'localhost';
}

export const config: Readonly<ConfigClass> = new ConfigClass();
