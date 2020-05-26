import {ParseProfile, Profile} from "../models/profile";

interface IAPIClient {
    GetProfile(handler: (profile: Profile) => void): void
}

const buildUrl = (host: string, path: string): string => {
    return host + path;
};

const NewAPI = (hostString: string): IAPIClient => {
    const u = new URL(hostString)
    const host = u.host
    return {
        GetProfile(handler: (profile: Profile) => void) {
            console.log("get profile")
            fetch(buildUrl(host, "/api/profile"), {
                method: "GET",
                redirect: "manual",
                headers: {
                    "Content-Type": "application/json"
                }
            }).then(r => {
                if (r.redirected) {
                    console.log(r)
                    window.location.href = r.url;
                } else if (r.status !== 200) {
                    console.log(r)
                    return
                }
                r.json().then(data => {
                    handler(ParseProfile(data))
                }
                ).catch(f => { console.log("failed to get/parse json" + f) })
            }).catch(f => {
                console.log("request failed: " + f);
            })
        }
    };
}

export { NewAPI }