import {ParseProfile, Profile} from "../models/profile";
import {MetricStatsRow, ParseMetricStats} from "../models/metric";

interface IAPIClient {
    GetProfile(handler: (profile: Profile) => void, errorHandler?: ErrorHandlerType): void
    GetMetricStats(bucket_size: number, handler: (stats: MetricStatsRow[]) => void, errorHandler?: ErrorHandlerType): void
}

const buildUrl = (host: string, path: string): string => {
    return host + path;
};

function maybeRedirectOrHandle<T>(response: Response, parser: (data: any) => T, handleData: DataHandlerType<T>, errorHandler: ErrorHandlerType): void {
    if (response.status === 401) {
        response.text().then(t => {
            window.location.href = t
        }).catch(f => errorHandler("failed to get redirect url: " + f))
    } else if (response.status !== 200) {
        errorHandler("status_code: " + response.status + " url: " + response.url)
    } else {
        response.json().then(data => {
            console.log(response, data)
            handleData(parser(data))
        }).catch(errorHandler)
    }
}

type DataHandlerType<T> = (data: T) => void

type ErrorHandlerType = (err: any) => void

const defaultErrorHandler = (err: any): void => {
    console.log(err)
}

const NewAPI = (hostString: string): IAPIClient => {
    const u = new URL(hostString)
    const host = u.host
    return {
        GetProfile(handler: (profile: Profile) => void, errorHandler: ErrorHandlerType = defaultErrorHandler) {
            fetch(buildUrl(host, "/api/profile"), {
                method: "GET",
                redirect: "follow",
                headers: {
                    "Content-Type": "application/json"
                }
            }).then(r => {
                maybeRedirectOrHandle(r, ParseProfile, handler, errorHandler)
            }).catch(f => {
                errorHandler("request failed: " + f)
            })
        },
        GetMetricStats(bucket_size: number, handler: (stats: MetricStatsRow[]) => void, errorHandler: ErrorHandlerType = defaultErrorHandler) {
            fetch(buildUrl(host, "/api/metric/stats"), {
                method: "POST",
                redirect: "follow",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    seconds: bucket_size,
                }),
            }).then(r => {
                maybeRedirectOrHandle(r, ParseMetricStats, handler, errorHandler)
            }).catch(f => {
                errorHandler("request failed: " + f)
            })
        }
    };
}

export { NewAPI }