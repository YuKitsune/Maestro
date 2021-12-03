import {LoaderFunction, MetaFunction, useLoaderData} from "remix";
import type {Thing} from "~/model/thing";
import CatalogueItem from "~/components/catalogueItem";
import HomeButton from "~/components/homeButton";
import {findBestThing} from "~/model/thing";

export let loader: LoaderFunction = async ({ params }) => {

    if (params.group === undefined) {
        throw new Error("Huston, we have a problem...")
    }

    const apiUrl = `${process.env.API_URL}/${params.group}`
    const res = await fetch(apiUrl)
    const json = await res.json()

    return json as Thing[];
};

export const meta: MetaFunction = ({data} : { data : Thing[] | undefined}) => {
    if (!data) {
        return {
            title: "Nothing found...",
            description: "🤦‍"
        };
    }

    const things = data as Thing[]
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
    let data = useLoaderData();
    let things = data as Thing[]

    return (
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            <CatalogueItem items={things} />
        </div>
    );
}
