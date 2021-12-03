import {findBestThing, Thing} from "~/model/thing";
import React, {useCallback} from "react";
import {Artist as ArtistModel} from "~/model/artist";
import {Album as AlbumModel} from "~/model/album";
import {Track as TrackModel} from "~/model/track";
import CatalogueItemLink from "~/components/catalogueItemLink";
import ArtistPreview from "~/components/artistPreview";
import AlbumPreview from "~/components/albumPreview";
import TrackPreview from "~/components/trackPreview";

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
    let bestThing = findBestThing(things);
    let preview: React.ReactNode;

    switch (bestThing!.ThingType) {
        case "artist":
            preview = <ArtistPreview artist={bestThing as ArtistModel} />
            break

        case "album":
            preview = <AlbumPreview album={bestThing as AlbumModel} />
            break

        case "track":
            preview = <TrackPreview track={bestThing as TrackModel} />
            break

        default:
            // Todo: Custom unknown component
            preview = <div>Unknown</div>
            break
    }

    return (
        <div className={"flex flex-col"}>

            {/* Preview */}
            <div className={"bg-gray-200 dark:bg-gray-700 rounded-lg p-2 mb-1"}>
                {preview}
            </div>

            {/* Links */}
            <div className={"flex flex-col gap-1"}>
                {things.map(t =>
                    <CatalogueItemLink key={t.Link} link={t.Link}>
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
                    </CatalogueItemLink>
                )}
            </div>
        </div>
    );
}

export default CatalogueItem;