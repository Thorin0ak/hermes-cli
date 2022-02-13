#!/usr/bin/env node

import chalk from 'chalk';
import inquirer from 'inquirer';
import gradient from 'gradient-string';
import ProgressBar from 'progress';
import { Command, Option } from 'commander';
import chalkAnimation from 'chalk-animation';
import figlet from 'figlet';
import {sleep} from './src/utils.js';
import {loadEnvironments} from './src/environments.js';
import {logger, stdoutTmpFileName, tmpFileName, configDir} from './src/initialization.js';
import {SseTest} from './src/test.js';

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';

const program = new Command();
program
    .name('mercure-test')
    .description('CLI to publish events to Mercure Hub')
    .version('0.1.0');
program
    .option('-o, --occurrences <number>', 'number of events to publish')
    .addOption(new Option('-o, --occurrences <number>', 'number of events to publish').argParser(parseInt))
    .option('-uri, --topic-uri <uri>', 'topic URI for subscription', 'sse://pxc.dev/123456/{term}')
    .addOption(new Option('-t, --type <data>', 'event type')
                    .choices(['mock', 'slots', 'visit-summaries'])
                    .default('random', 'mock payload'));

let environments, envName, envData;
const testConfiguration = {
    minWaitTime: 0,
    maxWaitTime: 2000,
    occurrences: 5,
    topicUri: 'sse://pxc.dev/123456/{term}'
};
let isDebugEnabled = true;

const strInitTool = 'Initializing the SSE testing tool...';
const strTestRunning = 'Sending data to Mercure Hub...';
const strChooseEnv = 'Which environment are you using?';

const initializeTool = async () => {
    const options = program.opts();
    if (options.topicUri) {
        testConfiguration.topicUri = options.topicUri;
    }
    if (options.occurrences) {
        testConfiguration.occurrences = options.occurrences;
    }
    console.clear();
    const rainbowTitle = chalkAnimation.rainbow(strInitTool);
    environments = loadEnvironments(configDir);
    await sleep(100); // just for fun  :)
    rainbowTitle.stop();
}

const selectEnvironment = async () => {
    const answers = await inquirer.prompt({
        name: 'envName',
        type: 'list',
        message: strChooseEnv,
        choices: Array.from(environments.keys())
    });

    envName = answers.envName;

    console.clear();
    const msg = `Environment: ${envName}`;
    envData = environments.get(envName);

    figlet(msg, (err, data) => {
        console.log(gradient.pastel.multiline(data));
        Object
            .keys(envData)
            .forEach(key => {
                if (envData[key].displayName) {
                    console.log(`${chalk.red(envData[key].displayName)}:  ${chalk.bgRed(envData[key].value)}`);
                }
            });
        console.log(`${chalk.red('TOPIC URI')}:  ${chalk.bgRed(testConfiguration.topicUri)}`);
        console.log('\n');
    });
}

const startTest = async () => {
    const progressBar = new ProgressBar(
        `${strTestRunning} [:bar] :percent | ETA: :eta seconds | :total/:current`,
        {stream: process.stdout, total: testConfiguration.occurrences}
    );

    SseTest.init(testConfiguration, envData);
    await SseTest.run(progressBar);
    if (progressBar.complete) {
        console.log('\ndone!\n');
    }
}

try {
    program.parse();
    await initializeTool();
    await selectEnvironment();
    await startTest();
} catch (error) {
    if (error.isTtyError) {
        console.log('Your console environment is not supported!');
    } else {
        logger.error(error);
        console.log(`An error occurred: ${error.message}`);
        console.log(`You can find more details here: ${tmpFileName}`);
        console.log(`Node stderr: ${stdoutTmpFileName}`);
    }
    process.exit(0);
}
