import http from 'k6/http';
import { check } from 'k6';

export function Check(serverUrl, headers, layer, base){
    const relationUrl = `${serverUrl}/relation`
    let res, payload;
    const namespace = "role", relation = "parent";
    const start = "0_0";
    let end = (layer).toString() + "_" + (Math.pow(base, layer)-1).toString();


    payload = {
        object_namespace: namespace,
        object_name: end,
        relation: relation,
        subject_namespace: namespace,
        subject_name: start,
        subject_relation: relation,
    };
    res = http.post(`${relationUrl}/check`, JSON.stringify(payload), {
        headers: headers, 
        timeout: '900s',
    });
    check(res, { 'Check': (r) => r.status ==  200 });
};