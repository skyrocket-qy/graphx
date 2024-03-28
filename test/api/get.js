

export function TestGetAPI(serverUrl, headers){
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
}