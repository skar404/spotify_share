import http from 'k6/http';

import {sleep} from 'k6';
import {Rate} from 'k6/metrics';

const myFailRate = new Rate('failed requests');


export let options = {
    vus: 50,
    duration: '10m',
    thresholds: {
        'failed requests': ['rate<0.1'], // threshold on a custom metric
        http_req_duration: ['p(95)<500'], // threshold on a standard metric
    },
};


export default function () {
    let res = http.get('http://localhost:1323/ping');
    myFailRate.add(res.status !== 200);

    sleep(1);
}
