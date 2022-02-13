import { mkdtemp, mkdir, copyFile } from 'fs/promises';
import { createWriteStream, existsSync } from 'fs';
import { sep, join } from 'path';
import { tmpdir, homedir } from 'os';
import winston from 'winston';

// handle config location: ~/.pxcdev/sse-cli/
const dir = `${homedir()}${sep}.pxcdev`
if (!existsSync(dir)) {
    await mkdir(dir);
    await mkdir(`${dir}${sep}sse-cli`);
}

export const configDir = `${dir}${sep}sse-cli`;

if (!existsSync(`${configDir}${sep}config.json`)) {
    await copyFile('sample-config.json', `${configDir}${sep}config.json`);
}

// temp dir and files for logging -- created with each run of the tool
const tmpFilePath = await mkdtemp(join(tmpdir(), 'sse-cli-'));
export const tmpFileName = `${tmpFilePath}${sep}winston.log`;
export const stdoutTmpFileName = `${tmpFilePath}${sep}node.log`;

// Handle stdout for Node warning, errors, etc...
const access = createWriteStream(stdoutTmpFileName);
process.stderr.write = access.write.bind(access);

export const logger = winston.createLogger({
    level: 'debug',
    format: winston.format.json(),
    defaultMeta: { service: 'sse-cli' },
    transports: [
        //
        // - Write all logs with importance level of `error` or less to `error.log`
        // - Write all logs with importance level of `info` or less to `combined.log`
        //
        // new winston.transports.File({ filename: 'error.log', level: 'error' }),
        new winston.transports.File({ filename: tmpFileName }),
    ],
});
