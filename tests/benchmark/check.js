import http from 'k6/http';
import { check } from 'k6';

export function Check(serverUrl, headers, layer, base){
    const relationUrl = `${serverUrl}/relation`
    const start = "0_0";
    let end = (layer).toString() + "_" + (Math.pow(base, layer)-1).toString();


    let payload = {
        object_namespace: "role",
        object_name: end,
        relation: "modify-permission",
        subject_namespace: "role",
        subject_name: start,
        subject_relation: "",
    };
    let res = http.post(`${relationUrl}/check`, JSON.stringify(payload), {
        headers: headers, 
        timeout: '900s',
    });
    check(res, { 'Check': (r) => r.status ==  200 });
};