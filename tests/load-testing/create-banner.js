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
    vus: 3,
    duration: '20s',
    gracefulStop: '3s',
};
export default function () {

    let tag = Math.floor(Math.random() * 4) + 1
    let feature = Math.floor(Math.random() * 15)
    // 10% очень важных пользователей
    let reallyImportantUser = Math.floor(Math.random() * 15)

    let newBanner : {
        tag_ids: [0],
        feature_id: 0,
        content: {
            title: "some_title",
            text: "some_text",
            url: "some_url"
        },
        is_active: true
    }

    const headers = { 'Content-Type': 'application/json', 'token': 'admin_token' };
    let res = http.post(`${env.backendUrl}/banner`, { headers })

    // 404 also is Ok because of no such banner
    http.setResponseCallback(http.expectedStatuses(200, 404))
}