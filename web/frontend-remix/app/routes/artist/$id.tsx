import {LoaderFunction, MetaFunction, useLoaderData} from "remix";
import MaestroApiClient from "~/maestroApiClient";
import {findBestThing} from "~/model/thing";
import {Artist} from "~/model/artist";
import HomeButton from "~/components/HomeButton";
import React from "react";
import ArtistPreview from "~/components/ArtistPreview";
import Card from "~/components/Card";

export let loader: LoaderFunction = async ({ params }) => {
    if (params.id === undefined) {
        throw new Error("Missing ID");
    }

    const client = new MaestroApiClient(process.env.API_URL as string, process.env.PUBLIC_API_URL as string);

    const artists = await client.getArtists(params.id);
    const services = await client.getServices();

    return {
        artists,
        services
    };
};

export const meta: MetaFunction = ({data}) => {
    const artists = data.artists as Artist[];
    let bestArtist = findBestThing<Artist>(artists);

    let title = bestArtist.Name;
    let image = bestArtist.ArtworkLink;
    let description = `Find ${title} on ${artists.length} streaming service${artists.length > 1 ? "s" : ""}!`;
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

export default function Artist() {
    let {artists, services} = useLoaderData();
    let bestArtist = findBestThing<Artist>(artists)

    return (
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            {artists === undefined && <h1 className={"text-lg"}>Sorry, couldn't find anything...</h1>}
            <Card renderPreview={() => <ArtistPreview artist={bestArtist} />} items={artists} services={services} />
        </div>
    );
}