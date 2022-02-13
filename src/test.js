import axios from 'axios';
import {getJwtToken, randomWait} from './utils.js';

class Test {
    axiosClient;
    testConfiguration;

    constructor(testConfiguration, envData) {
        this.testConfiguration = testConfiguration;
        const jwtToken = getJwtToken(testConfiguration.topicUri, envData.jwtSecretKey.value);
        this.axiosClient = axios.create({
            baseURL: envData.mercureHubUrl.value,
            headers: {
                Accept: 'application/json+fhir',
                Authorization: `Bearer ${jwtToken}`
            }
        });
    }

    async executeTest(progressBar) {
        const { occurrences, minWaitTime, maxWaitTime } = this.testConfiguration;
        for (let i = 0; i < occurrences; i++) {
            await randomWait(minWaitTime, maxWaitTime);
            await this.publish();
            progressBar.tick();
        }
    }


    async publish() {
        const params = new URLSearchParams({
            topic: this.testConfiguration.topicUri,
            data: JSON.stringify({foo: 'bar'}),
            type: 'scheduling_slots',
            private: 'on',
        });

        await this.axiosClient.post('mercure', params);
    }
}

let testInstance;

export const SseTest = {
    init: (topicUri, envData) => {
        if (!testInstance || !testInstance instanceof Test) {
            testInstance = new Test(topicUri, envData);
        }
    },
    run: async (progressBar) => await testInstance.executeTest(progressBar)
}
