import {LoaderFunction, MetaFunction, redirect} from "remix";
import Thing, {findBestThing, formatArtistNames} from "~/model/thing";
import MaestroApiClient from "~/maestroApiClient";
import {Artist} from "~/model/artist";
import {Album} from "~/model/album";
import {Track} from "~/model/track";
import Spinner from "~/components/spinner";
import {DefaultPort} from "~/defaults";

export let loader: LoaderFunction = async ({ params }) => {
    if (params.group === undefined) {
        throw new Error("Missing group ID")
    }

    const client = new MaestroApiClient(
        process.env.API_URL as string,
        `http://localhost:${process.env.PORT || DefaultPort}/api`);

    const id = params.group;

    const artistsRes = await client.tryGetArtists(id)
    const artists = artistsRes as Artist[];
    if (artists !== undefined && artists.length > 0) {
        return redirect(`/artist/${id}`)
    }

    const albumsRes = await client.tryGetAlbums(id)
    const albums = albumsRes as Album[];
    if (albums !== undefined && albums.length > 0) {
        return redirect(`/album/${id}`)
    }

    const tracksRes = await client.tryGetTracks(id)
    const tracks = tracksRes as Track[];
    if (tracks !== undefined && tracks.length > 0) {
        return redirect(`/track/${id}`)
    }

    throw new Error("Not sure what this ID is for...")
}

export const meta: MetaFunction = ({data}) => {
    const {type, things} = data;
    let bestThing = findBestThing(things);

    let title = getTitle(type, bestThing);
    let image = bestThing.ArtworkLink;
    let description = `Find ${title} on ${things.length} streaming service${things.length > 1 ? "s" : ""}!`;
    return {

        // Opengraph
        title: title,
        "og:title": title,
        "og:image": image,
        "og:site_name": "Maestro",
        "og:description": description,

        // Twitter
        "twitter:card": "summary",
        "twitter:title": title,
        "twitter:image": image,
        "twitter:description": description,
    };
};

function getTitle(type: string, thing: Thing) : any {
    switch (type) {
        case "artist":
        {
            return (thing as Artist).Name
        }

        case "album":
        {
            let album = (thing as Album)
            let artistNames = formatArtistNames(album.ArtistNames);
            return `${album.Name} by ${artistNames}`
        }

        case "track":
        {
            let track = (thing as Track)
            let artistNames = formatArtistNames(track.ArtistNames);
            return `${track.Name} by ${artistNames}`
        }
    }
}

export default function Group() {
    return <>
        <div className={"text-align-center"}>Looks like you've found an old link! Redirecting...</div>
        <Spinner className={"w-12 h-12"}/>
    </>
}
