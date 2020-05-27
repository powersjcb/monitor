export type MetricStatsRow = {
    Source: string,
    Name: string,
    TsBucket: number, // unix timestamp
    Avg: number,
    Max: number,
    Min: number,
}

const ParseMetricStatsRow = (data: any): MetricStatsRow => {
    return {
        Source: String(data.source),
        Name: String(data.name),
        TsBucket: Number(data.ts_bucket),
        Avg: Number(data.avg),
        Max: Number(data.max),
        Min: Number(data.min),
    }
};

export const ParseMetricStats = (data: any): Array<MetricStatsRow> => {
    console.log(data)
    const res: Array<MetricStatsRow> = [];
    if (!Array.isArray(data)) {
        return res;
    }
    for (let i = 0; i < data.length; i++) {
        if (data[i] !== undefined) {
            res.push(ParseMetricStatsRow(data[i]))
        }
    }
    return res;
};
