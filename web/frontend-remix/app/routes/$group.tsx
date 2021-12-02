import {Link, LoaderFunction, useLoaderData, useNavigate} from "remix";
import type {Thing} from "~/model/thing";
import {useCallback} from "react";
import CatalogueItem from "~/components/CatalogueItem";
import ArrowIcon, {ArrowDirection} from "~/components/icons/arrow";

export let loader: LoaderFunction = async ({ params }) => {

    if (params.group === undefined) {
        throw new Error("Huston, we have a problem...")
    }

    const apiUrl = `${process.env.API_URL}/${params.group}`
    const res = await fetch(apiUrl)
    const json = await res.json()

    return json as Thing[];
};

export default function Group() {
    let data = useLoaderData();
    let things = data as Thing[]

    const navigate = useNavigate()
    const goHome = useCallback(() => {
        navigate("/");
    }, [])

    return (
        <div className={"flex flex-col gap-2"}>

            <div className={"flex flex-initial"}>
                <Link to={"/"} className={"flex flex-row flex-initial content-center border-b border-opacity-0 hover:border-opacity-100"}>
                    <ArrowIcon direction={ArrowDirection.Left} className={"h-6 w-6 pr-2"} />
                    <div>Back to Maestro</div>
                </Link>
            </div>

            <CatalogueItem items={things} />
        </div>
    );
}
