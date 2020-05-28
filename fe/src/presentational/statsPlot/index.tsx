import React from "react";

import { VictoryAxis, VictoryChart, VictoryLabel, VictoryLine } from 'victory';
import {MetricStatsRow} from "../../lib/models/metric";

type PlotPoint = {
    x: number,
    y: number,
}

const mapStatsToPoints = (stats: MetricStatsRow[]): PlotPoint[] => {
    const res: PlotPoint[] = []
    for (let i=0; i < stats.length; i++) {
        res.push({
            x: stats[i].TsBucket,
            y: stats[i].Avg,
        })
    }
    return res
}

export const StatsPlot: React.FunctionComponent<{stats: MetricStatsRow[]}> = ({stats}) => {
    return (
        <VictoryChart >
            <VictoryLine
                data={mapStatsToPoints(stats)}
            />
            <VictoryAxis
                tickCount={3}
                tickFormat={(t) => {
                    const d = new Date(t * 1000)
                    console.log(t)
                    return d.toLocaleString("en-US")
                }}
                tickLabelComponent={<VictoryLabel angle={90} textAnchor={'start'} />}
            />
        </VictoryChart>
    );
}