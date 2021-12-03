
export type ThingType = "artist" | "album" | "track";

export type Thing = {
    Name: string
    ArtworkLink: string
    ThingType: ThingType
    GroupId: string
    Source: string
    Market: string
    Link: string
}

// Todo: Fix API so we don't need this, all results should be solid
export const findBestThing = (things: Thing[]): Thing => {
    let bestThing = things.find(t => t.ArtworkLink && t.ArtworkLink.length > 0)!;
    return bestThing;
}
