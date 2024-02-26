import http from 'k6/http';

const url = 'http://localhost:8080/hello';

export default function () {
    let data = { name: 'Tom' };

    const cb200 = http.expectedStatuses(200);
    // 传入 json
    http.post(url, JSON.stringify(data), {
        headers: { 'Content-Type': 'application/json' },
        responseCallback: cb200,
    });
}
