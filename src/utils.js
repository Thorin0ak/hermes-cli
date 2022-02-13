import jwt from 'jsonwebtoken';

export const sleep = (ms = 2000) => new Promise(resolve => setTimeout(resolve, ms));

export const getJwtToken = (topicUri, jwtSecretKey) => {
    const payload = {
        mercure: {publish: [topicUri]},
        sub: '123456',
        exp: new Date().getTime() / 1000 + 900000,
    }

    return jwt.sign(payload, jwtSecretKey);
}

export const randomWait = (a, b) => {
    const randomMsDelay = Math.floor(Math.random() * b) + a;
    return new Promise(resolve => setTimeout(resolve, randomMsDelay));
}
