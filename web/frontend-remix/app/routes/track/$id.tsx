import {LoaderFunction, MetaFunction, useLoaderData} from "remix";
import MaestroApiClient from "~/maestroApiClient";
import {findBestThing, formatArtistNames} from "~/model/thing";
import HomeButton from "~/components/homeButton";
import React from "react";
import Card from "~/components/card";
import {Track} from "~/model/track";
import TrackPreview from "~/components/trackPreview";
import {DefaultPort} from "~/defaults";

export let loader: LoaderFunction = async ({ params }) => {
    if (params.id === undefined) {
        throw new Error("Missing ID");
    }

    const client = new MaestroApiClient(
        process.env.API_URL as string,
        `http://localhost:${process.env.PORT || DefaultPort}/api`);

    const tracks = await client.getTracks(params.id);
    const services = await client.getServices();

    return {
        tracks,
        services
    };
};

export const meta: MetaFunction = ({data}) => {
    const tracks = data.tracks as Track[];
    let bestTrack = findBestThing<Track>(tracks)

    let title = `${bestTrack.Name} by ${formatArtistNames(bestTrack.ArtistNames)}`;
    let image = bestTrack.ArtworkLink;
    let description = `Find ${title} on ${tracks.length} streaming service${tracks.length > 1 ? "s" : ""}!`;
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
    let {tracks, services} = useLoaderData();
    let bestTrack = findBestThing<Track>(tracks);

    return (
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            {tracks === undefined && <h1 className={"text-lg"}>Sorry, couldn't find anything...</h1>}
            <Card renderPreview={() => <TrackPreview track={bestTrack} />} items={tracks} services={services} />
        </div>
    );
}