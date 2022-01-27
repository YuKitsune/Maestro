import {LoaderFunction, MetaFunction, useLoaderData} from "remix";
import type {Thing} from "~/model/thing";
import CatalogueItem from "~/components/catalogueItem";
import HomeButton from "~/components/homeButton";
import {findBestThing} from "~/model/thing";
import {Service} from "~/model/service";
import MaestroApiClient from "~/maestroApiClient";

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

// @ts-ignore
export const meta: MetaFunction = ({data}) => {
    const things = data.things;
    let bestThing = findBestThing(things)

    let title = bestThing.Name;
    let image = bestThing.ArtworkLink;
    return {
        title: title,
        "og:title": title,
        "og:image": image,
        "og:site_name": "Maestro",
        "og:description": "Share music regardless of streaming service!",
    };
};

export default function Group() {
    let {things, services} = useLoaderData<GroupData>();

    return (
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            <CatalogueItem things={things} services={services} />
        </div>
    );
}
