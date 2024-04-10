import http from 'k6/http';
import { env } from './config.js';

// export const options = {
//     scenarios: {
//         // my_scenario1: {
//         //     executor: 'constant-arrival-rate',
//         //     duration: '5s', // total duration
//         //     preAllocatedVUs: 100, // to allocate runtime resources     preAll
//         //     gracefulStop: '5s',
//         //     rate: 1000, // number of constant iterations given `timeUnit`
//         //     timeUnit: '1s',
//         // },
//         contacts: {
//             executor: 'constant-vus',
//             vus: 1000,
//             duration: '30s',
//             gracefulStop: '5s',
//         },
//     },
// };
export const options = {
    vus: 10,
    duration: '20s',
    gracefulStop: '3s',
};
export default function () {

    let tag = Math.floor(Math.random() * 4)
    let feature = Math.floor(Math.random() * 15)
    // 10% очень важных пользователей
    let reallyImportantUser = Math.floor(Math.random() * 15)
    const headers = { 'Content-Type': 'application/json', 'token': 'user_token' };
    let res = http.get(`${env.backendUrl}/user_banner?tag_id=${tag}&feature_id=${feature}&use_last_revision=${reallyImportantUser > 8}`, { headers })

    // 404 also is Ok because of no such banner. 403 - если токен юзера и баннер не активен
    http.setResponseCallback(http.expectedStatuses(200, 403, 404))
}

