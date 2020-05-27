import React, {useEffect, useState} from "react";

import {ProfileHeader} from "../../presentational/profileHeader";
import {Profile} from "../../lib/models/profile";
import {NewAPI} from "../../lib/api";
import {MetricStatsRow, ParseMetricStats} from "../../lib/models/metric";

const LandingPage: React.FunctionComponent = () => {
    const [profile, setProfile] = useState<Profile>({AccountID: 0, ApiKey: ""})
    const [stats, setStats] = useState<Array<MetricStatsRow>>([])
    useEffect(() => {
        NewAPI("localhost:3000").GetProfile(setProfile)
        NewAPI("localhost:3000").GetMetricStats(setStats)
    }, [])
    return (
        <div>
            <ProfileHeader profile={profile} />
            {stats}
        </div>
    )
}

export { LandingPage }