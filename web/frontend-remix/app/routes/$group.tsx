import {LoaderFunction, MetaFunction, useLoaderData} from "remix";
import type {Thing} from "~/model/thing";
import {findBestThing, formatArtistNames} from "~/model/thing";
import CatalogueItem from "~/components/catalogueItem";
import HomeButton from "~/components/homeButton";
import {Service} from "~/model/service";
import MaestroApiClient from "~/maestroApiClient";
import {Artist} from "~/model/artist";
import {Album} from "~/model/album";
import {Track} from "~/model/track";

type GroupData = {
    things: Thing[];
    services: Service[];
}

export let loader: LoaderFunction = async ({ params }) => {
    if (params.group === undefined) {
        throw new Error("Huston, we have a problem...")
    }

    const client = new MaestroApiClient(process.env.API_URL as string, process.env.PUBLIC_API_URL as string)

    const things = await client.getGroup(params.group)
    const services = await client.getServices()

    return {
        things,
        services
    };
};

export const meta: MetaFunction = ({data}) => {
    const things = data.things;
    let bestThing = findBestThing(things)

    let title = getTitleForThing(bestThing);
    let image = bestThing.ArtworkLink;
    let description = `Listen to ${title} on ${things.length} streaming service${things.length > 1 ? "s" : ""}!`;
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

function getTitleForThing(thing: Thing) : any {
    switch (thing.ThingType) {
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
    let {things, services} = useLoaderData<GroupData>();

    return (
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            <CatalogueItem things={things} services={services} />
        </div>
    );
}
