import React, {useEffect, useState} from "react";

import {ProfileHeader} from "../../presentational/profileHeader";
import {Profile} from "../../lib/models/profile";
import {NewAPI} from "../../lib/api";

const LandingPage: React.FunctionComponent = () => {
    const [profile, setProfile] = useState({AccountID: 0, ApiKey: ""})
    useEffect(() => {
        const handleProfileResp = (profile: Profile) => {
            setProfile(profile);
        };
        NewAPI("localhost:3000").GetProfile(handleProfileResp)
    }, [])
    return (
        <div>
            <ProfileHeader profile={profile} />
        </div>
    )
}

export { LandingPage }