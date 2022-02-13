import { readFileSync } from 'fs';
import { sep } from 'path';

export const loadEnvironments = (configDir) => {
    const configJson = readFileSync(`${configDir}${sep}config.json`, { encoding: 'utf-8' });
    const config = JSON.parse(configJson);

    const environments = new Map();
    config.Environments.forEach(env => environments.set(env.name, env.config));

    return environments;
}
