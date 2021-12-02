import {useNavigate} from "remix";
import {Thing} from "~/model/thing";
import {useCallback} from "react";
import {Artist as ArtistModel} from "~/model/artist";
import {Album as AlbumModel} from "~/model/album";
import {Track as TrackModel} from "~/model/track";
import Link from "~/components/link";
import Artist from "~/components/artist";
import Album from "~/components/album";
import Track from "~/components/track";

type CatalogueItemProps = {
    items: Thing[]
}

const CatalogueItem = (props: CatalogueItemProps) => {
    let things = props.items;

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

    }, [props.items])

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

    return (
        <div className={"flex flex-col"}>

            {/* Preview */}
            <div className={"mb-1"}>
                {preview}
            </div>

            {/* Links */}
            <div className={"flex flex-col gap-1"}>
                {things.map(t =>
                    <Link key={t.Link} link={t.Link}>
                        <div className={"flex flex-row content-center gap-1"}>

                            {/* Logo */}
                            <img src={getSourceIconLink(t.Source)} alt={t.Source} className={"row-span-2 flex-shrink h-12"}/>

                            {/* Source name and link */}
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

export default CatalogueItem;