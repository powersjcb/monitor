export type Profile = {
    AccountID: number,
    ApiKey: string,
}

const ParseProfile = (data: any): Profile => {
    return {
        AccountID: parseInt(data.id),
        ApiKey: String(data.api_key),
    }
}

export { ParseProfile }