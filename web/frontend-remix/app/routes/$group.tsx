import {LoaderFunction, redirect, useLoaderData, useNavigate} from "remix";
import type {Thing} from "~/model/thing";
import type {Artist as ArtistModel} from "~/model/artist";
import type {Album as AlbumModel} from "~/model/album";
import type {Track as TrackModel} from "~/model/track";
import Artist from "~/components/artist";
import Album from "~/components/album";
import Track from "~/components/track";
import Link from "~/components/link";
import {useCallback} from "react";

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

    // Todo: Blegh... Move these to a CDN or something, then have the API return a link along with them
    const getSourceIconLink = useCallback((sourceName: string) => {
        switch (sourceName) {
            case "Apple Music":
                return "/images/Apple Music.png";

            case "Spotify":
                return "/images/Spotify.png";

            case "Deezer":
                return "/images/Deezer.png";

            // Todo: question mark instead
            default:
                return ""
        }

    }, [data])

    if (things == undefined || things.length == 0) {
        return <h1 className={"text-lg"}>Sorry, couldn't find anything...</h1>
    }

    // Find the most appropriate thing
    let bestThing = things.find(t => t.ArtworkLink && t.ArtworkLink.length > 0);
    let preview: React.ReactNode;

    switch (bestThing!.ThingType) {
        case "artist":
            preview = <Artist artist={bestThing as ArtistModel} />
            break

        case "album":
            preview = <Album album={bestThing as AlbumModel} />
            break

        case "track":
            preview = <Track track={bestThing as TrackModel} />
            break

        default:
            // Todo: Custom unknown component
            preview = <div>Unknown</div>
            break
    }

    const navigate = useNavigate()
    const goHome = useCallback(() => {
        navigate("/");
    }, [])

    return (
        <div className={"flex flex-col gap-1"}>

            {/* Todo: Scuffed... A link with a back icon would be okay */}
            <div className={"bg-gray-200 dark:bg-gray-700 rounded-lg p-2 cursor-pointer"} onClick={goHome}>
                Home
            </div>

            {preview}

            {/* Links */}
            <div className={"flex flex-col gap-1 p-2"}>
                {things.map(t =>
                        <Link key={t.Link} link={t.Link}>
                            <div className={"flex flex-row content-center gap-1"}>

                                <img src={getSourceIconLink(t.Source)} alt={t.Source} className={"row-span-2 flex-shrink h-12"}/>

                                <div className={"grid"}>
                                    <div className={"font-bold"}>
                                        {t.Source}
                                    </div>
                                    <div className={"text-underline text-blue-400 truncate"}>
                                        {t.Link}
                                    </div>
                                </div>
                            </div>
                        </Link>
                )}
            </div>
        </div>
    );
}
