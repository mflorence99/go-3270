export class ConfigClass {
  port = '3000';
}

export const config: Readonly<ConfigClass> = new ConfigClass();
