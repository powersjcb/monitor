

import React from "react";

import {Profile} from "../../lib/models/profile";

type ProfileHeaderProps = {
    profile: Profile,
}

export const ProfileHeader: React.FC<ProfileHeaderProps> = ({profile}) => {
    return <div>accountID: {profile.AccountID}, apiKey: {profile.ApiKey}</div>
}
