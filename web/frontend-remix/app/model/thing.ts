
export default interface Thing {
    Name: string;
    Link: string;
    ArtworkLink: string;
    Source: string
}

// Todo: Fix API so we don't need this, all results should be solid
export function findBestThing<T extends Thing>(things: T[]): T {
    return things.find(t => t.ArtworkLink && t.ArtworkLink.length > 0)!;
}

export function formatArtistNames(names: string[]): string {
    if (names.length == 1) {
        return names[0];
    }

    let formattedNames = "";
    for (let i = 0; i < names.length; i++) {

        if (i == names.length - 1) {
            formattedNames += " and "
        } else if (i > 0) {
            formattedNames += ", "
        }

        formattedNames += names[i];
    }

    return formattedNames;
}
