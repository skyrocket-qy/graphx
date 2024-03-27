import http from 'k6/http';
import { check } from 'k6';


export function TestRelationAPI(serverUrl, headers){
    const relationUrl = `${serverUrl}/relation`
    let res, payload;

    res = http.get(`${relationUrl}?`, null, {headers:headers});
    check(res, { 'Get': (r) => r.status == 200 });

    payload = {
        object_namespace: "test_file",
        object_name: "1",
        relation: "read",
        subject_namespace: "test_file",
        subject_name: "1",
        subject_relation: "write",
    };
    res = http.post(`${relationUrl}`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Create': (r) => r.status == 200 });

    payload = {
        subject: {
            namespace: "test_file",
            name: "1",
            relation: "write",
        },
        object: {
            namespace: "test_file",
            name: "1",
            relation: "read",
        },
    };
    res = http.post(`${relationUrl}/check`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Check': (r) => r.status ==  200 });

    payload = {
        subject: {
            namespace: "test_file",
            name: "1",
            relation: "write",
        },
        object: {
            namespace: "test_file",
            name: "1",
            relation: "read",
        },
    };
    res = http.post(`${relationUrl}/get-shortest-path`, JSON.stringify(payload), {headers:headers});
    check(res, { 'GetShortestPath': (r) => r.status ==  200 });

    payload = {
        subject: {
            namespace: "test_file",
            name: "1",
            relation: "write",
        },
        object: {
            namespace: "test_file",
            name: "1",
            relation: "read",
        },
    };
    res = http.post(`${relationUrl}/get-all-paths`, JSON.stringify(payload), {headers:headers});
    check(res, { 'GetAllPaths': (r) => r.status ==  200 });

    payload = {
        subject: {
            namespace: "test_file",
            name: "1",
            relation: "write",
        },
    };
    res = http.post(`${relationUrl}/get-all-object-relations`, JSON.stringify(payload), {headers:headers});
    check(res, { 'GetAllObjectRelations': (r) => r.status ==  200 });

    payload = {
        object: {
            namespace: "test_file",
            name: "1",
            relation: "read",
        },
    };
    res = http.post(`${relationUrl}/get-all-subject-relations`, JSON.stringify(payload), {headers:headers});
    check(res, { 'GetAllSubjectRelations': (r) => r.status ==  200 });

    res = http.post(`${relationUrl}/get-all-namespaces`, null, {headers:headers});
    check(res, { 'GetAllNamespaces': (r) => r.status == 200 });

    payload = {
        object_namespace: "test_file",
        object_name: "1",
        relation: "read",
        subject_namespace: "test_file",
        subject_name: "1",
        subject_relation: "write",
    };
    res = http.del(`${relationUrl}`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Delete': (r) => r.status == 200 });

    res = http.post(`${relationUrl}/clear-all-relations`, null, {headers:headers});
    check(res, { 'ClearAllRelations': (r) => r.status == 200 });
}