import {LoaderFunction, MetaFunction, useLoaderData} from "remix";
import type {Thing} from "~/model/thing";
import CatalogueItem from "~/components/catalogueItem";
import HomeButton from "~/components/homeButton";
import {findBestThing} from "~/model/thing";
import {Service} from "~/model/service";

type GroupData = {
    things: Thing[];
    services: Service[];
}

const loadThings = async (groupId: string): Promise<Thing[]> => {
    const apiUrl = `${process.env.API_URL}/${groupId}`
    const res = await fetch(apiUrl)
    const json = await res.json()
    return json as Thing[];
}

const loadServices = async (): Promise<Service[]> => {
    const apiUrl = `${process.env.API_URL}/services`
    const res = await fetch(apiUrl)
    const json = await res.json()
    return json as Service[];
}

export let loader: LoaderFunction = async ({ params }) => {
    if (params.group === undefined) {
        throw new Error("Huston, we have a problem...")
    }

    const things = await loadThings(params.group)
    const services = await loadServices()

    return {
        things,
        services
    };
};

// @ts-ignore
export const meta: MetaFunction = ({things} : GroupData) => {
    if (!things) {
        return {
            title: "Nothing found...",
            description: "🤦‍"
        };
    }

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
