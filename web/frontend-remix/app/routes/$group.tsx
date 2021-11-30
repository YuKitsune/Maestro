import {LoaderFunction, useLoaderData} from "remix";
import type {Thing} from "~/model/thing";
import type {Artist as ArtistModel} from "~/model/artist";
import type {Album as AlbumModel} from "~/model/album";
import type {Track as TrackModel} from "~/model/track";
import Artist from "~/components/artist";
import Album from "~/components/album";
import Track from "~/components/track";

export let loader: LoaderFunction = async ({ params }) => {

    if (params.group === undefined) {
        throw new Error("Huston, we have a problem...")
    }

    const apiUrl = `${process.env.API_URL}/${params.group}`
    const res = await fetch(apiUrl)
    const json = await res.json()

    return json as Thing[];
};

export default function Links() {
    let data = useLoaderData();
    let things = data as Thing[]

    return (
        <div className={"flex flex-col gap-1"}>
            {things.map(t => {
                switch (t.ThingType) {
                    case "artist":
                        return <Artist key={t.Link} artist={t as ArtistModel} />
                    case "album":
                        return <Album key={t.Link} album={t as AlbumModel} />
                    case "track":
                        return <Track key={t.Link} track={t as TrackModel} />
                    default:
                        // Todo: Custom unknown component
                        return <div key={t.Link}>Unknown</div>
                }
            })}
        </div>
    );
}
