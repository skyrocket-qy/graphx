import http from 'k6/http';
import { check } from 'k6';


export function TestRelationAPI(serverUrl, headers){
    const relationUrl = `${serverUrl}/relation`
    let res, payload;
    let edge = {
        "obj-ns": "role",
        "obj-name": "rd",
        "obj-rel": "parent",
        "sbj-ns": "role",
        "sbj-name": "rd-director",
    };

    res = http.get(relationUrl +
                    "?obj-ns=" + edge["obj-ns"] +
                    "&obj-name=" + edge["obj-name"] +
                    "&obj-rel=" + edge["obj-rel"] +
                    "&sbj-ns=" + edge["sbj-ns"] +
                    "&sbj-name=" + edge["sbj-name"],
                    null, {headers:headers});
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
        object_namespace: "test_file",
        object_name: "1",
        relation: "read",
        subject_namespace: "test_file",
        subject_name: "1",
        subject_relation: "write",
    };
    res = http.del(`${relationUrl}`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Delete': (r) => r.status == 200 });

    res = http.del(`${relationUrl}/all`, JSON.stringify(payload), {headers:headers});
    check(res, { 'ClearAll': (r) => r.status == 200 });

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
    check(res, { 'CheckAuth': (r) => r.status ==  200 });

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
    res = http.post(`${relationUrl}/obj-auths`, JSON.stringify(payload), {headers:headers});
    check(res, { 'GetObjAuths': (r) => r.status ==  200 });

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
    res = http.post(`${relationUrl}/sbj-who-has-auth`, JSON.stringify(payload), {headers:headers});
    check(res, { 'GetSbjsWhoHasAuth': (r) => r.status ==  200 });

    payload = {
        subject: {
            namespace: "test_file",
            name: "1",
            relation: "write",
        },
    };
    res = http.post(`${relationUrl}/get-tree`, JSON.stringify(payload), {headers:headers});
    check(res, { 'GetTree': (r) => r.status ==  200 });

    res = http.post(`${relationUrl}/see-tree`, null, {headers:headers});
    check(res, { 'SeeTree': (r) => r.status == 200 });
}