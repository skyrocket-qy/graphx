import http from 'k6/http';
import { check } from 'k6';
  
export function checkScenario(serverUrl, headers){
    let payload = {
      user_name: "Jimmy",
      role_name: "rd-director",
    };
    // Add relation
    let res = http.post(`${serverUrl}/user/add-role`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) Jimmy is a member of RD-Director': (r) => r.status == 200 } );

    payload = {
      user_name: "Tasha",
      role_name: "rd",
    };
    res = http.post(`${serverUrl}/user/add-role`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) Tasha is a member of RD': (r) => r.status == 200 });

    payload = {
      user_name: "Ivy",
      role_name: "sales",
    };
    res = http.post(`${serverUrl}/user/add-role`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) Ivy is a member of Sales': (r) => r.status == 200 });

    payload = {
      child_role_name: "rd",
      parent_role_name: "rd-director",
    };
    res = http.post(`${serverUrl}/role/add-parent`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) RD-Director is a parent of RD': (r) => r.status == 200 });

    // Add relation
    payload = {
      user_name: "Ivy",
      relation: "read",
      object_namespace: "profile",
      object_name: "Ivy",
    };
    res = http.post(`${serverUrl}/user/add-relation`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) Ivy has read access to her profile': (r) => r.status == 200 });

    payload = {
      user_name: "Heidi",
      relation: "read",
      object_namespace: "profile",
      object_name: "Heidi",
    };
    res = http.post(`${serverUrl}/user/add-relation`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) Heidi has read access to her profile': (r) => r.status == 200 });

    payload = {
      role_name: "rd-director",
      relation: "*",
      object_namespace: "source-code",
      object_name: "*",
    };
    res = http.post(`${serverUrl}/role/add-relation`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) RD-Director has full access to all source code': (r) => r.status == 200 });

    payload = {
      role_name: "rd",
      relation: "write",
      object_namespace: "source-code",
      object_name: "*",
    };
    res = http.post(`${serverUrl}/role/add-relation`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) RD has write access to all source code': (r) => r.status == 200 });

    payload = {
      role_name: "rd-director",
      relation: "read",
      object_namespace: "sales-data",
      object_name: "*",
    };
    res = http.post(`${serverUrl}/role/add-relation`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) RD-Director has read access to all sales data': (r) => r.status == 200 });

    payload = {
      role_name: "sales",
      relation: "*",
      object_namespace: "sales-data",
      object_name: "*",
    };
    res = http.post(`${serverUrl}/role/add-relation`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) Sales has full access to all sales data': (r) => r.status == 200 });

    payload = {
      role_name: "rd",
      relation: "write",
      object_namespace: "sales-data",
      object_name: "1",
    };
    res = http.post(`${serverUrl}/role/add-relation`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Relation) RD has write access to sales-data:1': (r) => r.status == 200 });

    // Link relation
    // payload = {
    //   object_namespace: "role",
    //   object_name: "rd",
    //   relation: "modify-relation",
    //   subject_namespace: "role",
    //   subject_name: "rd",
    //   subject_relation: "parent",
    // };
    // res = http.post(`${serverUrl}/relation/link`, JSON.stringify(payload), {headers:headers});
    // check(res, { '(Add Link Relation) RD-Director can modify permissions of RD': (r) => r.status == 200 });

    payload = {
      object_namespace: "source-code",
      object_name: "*",
      relation: "read",
      subject_namespace: "source-code",
      subject_name: "*",
      subject_relation: "write",
    };
    res = http.post(`${serverUrl}/relation/link`, JSON.stringify(payload), {headers:headers});
    check(res, { '(Add Link Relation) Source Code read access implies write access': (r) => r.status == 200 });


    // Test 
    payload = {
      object_namespace: "profile",
      object_name: "Ivy",
      relation: "read",
      user_name: "Ivy"
    };
    res = http.post(`${serverUrl}/user/check`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Ivy has read to profile:Ivy': (r) => r.status == 200 && JSON.parse(r.body).result === true});

    payload = {
      object_namespace: "source-code",
      object_name: "2",
      relation: "delete",
      user_name: "Jimmy"
    };
    res = http.post(`${serverUrl}/user/check`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Jimmy has delete to source-code:2': (r) => r.status == 200 && JSON.parse(r.body).result === true});

    payload = {
      object_namespace: "role",
      object_name: "rd",
      relation: "modify-relation",
      user_name: "Jimmy"
    };
    res = http.post(`${serverUrl}/user/check`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Jimmy has modify-relation to role:rd': (r) => r.status == 200 && JSON.parse(r.body).result === true});

    payload = {
      object_namespace: "sales-data",
      object_name: "1",
      relation: "write",
      user_name: "Jimmy"
    };
    res = http.post(`${serverUrl}/user/check`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Jimmy has write to sales-data:1': (r) => r.status == 200 && JSON.parse(r.body).result === true});

    payload = {
      object_namespace: "source-code",
      object_name: "1",
      relation: "read",
      user_name: "Tasha"
    };
    res = http.post(`${serverUrl}/user/check`, JSON.stringify(payload), {headers:headers});
    check(res, { 'Jimmy has write to sales-data:1': (r) => r.status == 200 && JSON.parse(r.body).result === true});

    res = http.del(`${serverUrl}/relation/`, null, {headers:headers});
    check(res, { 'ClearAllRelations: status == 200': (r) => r.status == 200 });
}