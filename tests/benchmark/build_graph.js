import http from 'k6/http';
import { check } from 'k6';

export function BuildGraph(serverUrl, headers, layer, base){
    const roleUrl = `${serverUrl}/role`
    let payload, res;
    
    let curLayer = 1;
    while (curLayer <= layer){
        const count = Math.pow(base, curLayer);

        for (let i = 0; i < count; i++){
            payload = {
                parent_role_name: (curLayer-1).toString() + "_" + Math.floor(i / base).toString(),
                child_role_name: curLayer.toString() + "_" + i.toString(),
            };
            res = http.post(`${roleUrl}/add-parent`, JSON.stringify(payload), {headers: headers});
            check(res, { 'AddParent: status == 200': (r) => r.status == 200 });
        };

        curLayer += 1;
    };
};