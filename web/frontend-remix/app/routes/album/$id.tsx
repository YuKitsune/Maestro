import {LoaderFunction, MetaFunction, useLoaderData} from "remix";
import MaestroApiClient from "~/maestroApiClient";
import {findBestThing, formatArtistNames} from "~/model/thing";
import HomeButton from "~/components/homeButton";
import React from "react";
import Card from "~/components/card";
import {Album} from "~/model/album";
import AlbumPreview from "~/components/albumPreview";

export let loader: LoaderFunction = async ({ params }) => {
    if (params.id === undefined) {
        throw new Error("Missing ID");
    }

    const client = new MaestroApiClient(process.env.API_URL as string, process.env.PUBLIC_API_URL as string);

    const albums = await client.getAlbums(params.id);
    const services = await client.getServices();

    return {
        albums,
        services
    };
};

export const meta: MetaFunction = ({data}) => {

    console.log("Meta", data);

    const albums = data.albums as Album[];
    let bestAlbum = findBestThing<Album>(albums)

    let title = `${bestAlbum.Name} by ${formatArtistNames(bestAlbum.ArtistNames)}`;
    let image = bestAlbum.ArtworkLink;
    let description = `Find ${title} on ${albums.length} streaming service${albums.length > 1 ? "s" : ""}!`;
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

export default function Album() {
    let {albums, services} = useLoaderData();
    let bestAlbum = findBestThing<Album>(albums)

    return (
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            {albums === undefined && <h1 className={"text-lg"}>Sorry, couldn't find anything...</h1>}
            <Card renderPreview={() => <AlbumPreview album={bestAlbum} />} items={albums} services={services} />
        </div>
    );
}