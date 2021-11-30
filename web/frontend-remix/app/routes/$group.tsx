import {LoaderFunction, useLoaderData} from "remix";

export let loader: LoaderFunction = async ({ params }) => {

    if (params.link === undefined) {
        throw new Error("Huston, we have a problem...")
    }

    const encodedLink = encodeURIComponent(params.link)

    const apiUrl = `http://localhost:8182/link?link=${encodedLink}`
    const res = await fetch(apiUrl)
    const json = await res.json()

    return json
};

export default function Links() {
    let slug = useLoaderData();
    return (
        <div>
            <h1>Some Post: {slug}</h1>
        </div>
    );
}
