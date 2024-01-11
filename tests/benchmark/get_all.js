import http from 'k6/http';
import { check } from 'k6';


export const options = {
    vus: 100,
    duration: '180s',
}

function generateRandomData(i) {
    return {
        relation: {
            object_namespace: i,
            object_name: i,
            relation: i,
            subject_namespace: i,
            subject_name: i,
            subject_relation: i,
        },
        exist_ok: false,
    };
}

export function setup() {
    const headers = {
        'Content-Type': 'application/json',
    };
    const res = http.post(`http://localhost:8080/relation/clear-all-relations`, null, {headers:headers});
    check(res, { 'ClearAllRelations': (r) => r.status == 200 });

    // Clear all relations
    const clearRes = http.post(`http://localhost:8080/relation/clear-all-relations`, null, { headers: headers });
    check(clearRes, { 'ClearAllRelations': (r) => r.status == 200 });

    // Create new relations
    for (let i = 0; i < 5000; i++) {
        const randomData = generateRandomData(i.toString());
        const createRes = http.post(`http://localhost:8080/relation`, JSON.stringify(randomData), { headers: headers });
        check(createRes, { 'Create request was successful': (r) => r.status === 200 });
    }
}

export function testGetAll() {
    const headers = {
        'Content-Type': 'application/json',
    };

    for (let i = 0; i < 1000; i++) {
        const response = http.get(`http://localhost:8080/relation`, { headers: headers });
        check(response, { 'Query request was successful': (r) => r.status === 200 });
    }
}

export default function () {
    testGetAll();
}
