import {findBestThing, Thing} from "~/model/thing";
import React from "react";
import {Artist as ArtistModel} from "~/model/artist";
import {Album as AlbumModel} from "~/model/album";
import {Track as TrackModel} from "~/model/track";
import CatalogueItemLink from "~/components/catalogueItemLink";
import ArtistPreview from "~/components/artistPreview";
import AlbumPreview from "~/components/albumPreview";
import TrackPreview from "~/components/trackPreview";
import {Service} from "~/model/service";

type CatalogueItemProps = {
    things: Thing[]
    services: Service[]
}

const CatalogueItem = ({ things, services }: CatalogueItemProps) => {

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
                    <CatalogueItemLink key={t.Source} thing={t} service={services.find(s => s.Key === t.Source)}/>
                )}
            </div>
        </div>
    );
}

export default CatalogueItem;
