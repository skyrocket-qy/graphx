import http from 'k6/http';
import { check } from 'k6';
  
export function TestCycle(serverUrl, headers){
    const roleUrl = `${serverUrl}/role`
    const relation = "parent";
    const roleName1 = "Jimmy";
    const roleName2 = "Tasha"
    let payload;
    let res;

    payload = {
        child_role_name: roleName2,
        parent_role_name: roleName1,
    };
    res = http.post(`${roleUrl}/add-parent`, JSON.stringify(payload), {headers:headers});
    check(res, { 'AddParent: status == 200': (r) => r.status == 200 });

    payload = {
        child_role_name: roleName1,
        parent_role_name: roleName2,
    };
    res = http.post(`${roleUrl}/add-parent`, JSON.stringify(payload), {headers:headers});
    check(res, { 'AddParent reverse: status != 200': (r) => r.status != 200 });

    res = http.del(`${serverUrl}/relation/`, null, {headers:headers});
    check(res, { 'ClearAllRelations: status == 200': (r) => r.status == 200 });
}