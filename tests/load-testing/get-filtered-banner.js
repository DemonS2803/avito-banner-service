import http from 'k6/http';
import { env } from './config.js';

export const options = {
    vus: 5,
    duration: '20s',
    gracefulStop: '3s',
};

export default function () {
    let tag = Math.floor(Math.random() * 4)
    let feature = Math.floor(Math.random() * 15)
    const headers = { 'Content-Type': 'application/json', 'token': 'admin_token' };
    let test = Math.floor(Math.random() * 3)
    let res
    switch (test) {
        case 0:
            res = http.get(`http://localhost:8080/banner?tag_id=${tag}`, { headers })
            break
        case 1:
            res = http.get(`http://localhost:8080/banner?feature_id=${feature}`, { headers })
            break
        case 2:
            res = http.get(`http://localhost:8080/banner?tag_id=${tag}&feature_id=${feature}`, { headers })
            break
    }

    // 404 also is Ok because of no such banner
    http.setResponseCallback(http.expectedStatuses(200))
}