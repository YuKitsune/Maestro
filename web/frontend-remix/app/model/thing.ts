
export type ThingType = "artist" | "album" | "track";

export type Thing = {
    ArtworkLink: string
    ThingType: ThingType
    GroupId: string
    Source: string
    Market: string
    Link: string
}
