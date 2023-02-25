import type {ActionFunction, MetaFunction} from "@remix-run/server-runtime";
import {Form} from "@remix-run/react";
import MaestroApiClient from "~/maestroApiClient";
import {Artist} from "~/model/artist";
import {Album} from "~/model/album";
import {Track} from "~/model/track";
import {ClipboardEvent, useRef} from "react";
import PlayIcon from "~/components/icons/playIcon";
import {redirect} from "@remix-run/router";


// https://remix.run/api/conventions#meta
export let meta: MetaFunction = () => {
  return {
    title: "Maestro",
    description: "Share music across a variety of streaming services."
  };
};

export const action: ActionFunction = async ({request}) => {
    const client = new MaestroApiClient(process.env.API_URL as string, process.env.PUBLIC_API_URL as string)

    const formData = await request.formData();
    const link = formData.get("link") as string;
    const response = await client.searchFromLink(link);

    switch (response.Type) {
        case "artist":
            const artists = response.Items as Artist[];
            return redirect(`/artist/${artists[0].ArtistId}`)

        case "album":
            const albums = response.Items as Album[];
            return redirect(`/album/${albums[0].AlbumId}`)

        case "track":
            const tracks = response.Items as Track[];
            return redirect(`/track/${tracks[0].Isrc}`)

        default:
            return redirect(`/nothing`);
    }
}

// https://remix.run/guides/routing#index-routes
export default function Index() {
    const formRef = useRef<HTMLFormElement>(null);

    const onPaste = (event: ClipboardEvent) => {
        if (!event.clipboardData)
            return;

        setTimeout(async () => {
            formRef?.current?.submit();
        }, 500);
    }

    return (
        <div className={"flex flex-col gap-2 align-center items-center"}>
            <h1 className={"text-4xl text-align-center"}>Maestro</h1>
            <div className={"text-center"}>
                Share links to artists, albums, or tracks across a variety of different streaming services.
            </div>

            <Form method={"post"} className={"w-full flex flex-row gap-2"} ref={formRef}>
                <input
                    className={"w-full rounded-lg border-2 border-blue-100 dark:border-blue-900 dark:bg-black px-1 focus-within:border-blue-500 outline-none"}
                    type="text"
                    placeholder={"Paste URL"}
                    name={"link"}
                    onPaste={onPaste}
                />
                <button type="submit">
                    <PlayIcon className={"h-6 w-6 fill-green-300 hover:fill-green-400 dark:fill-green-800 dark:hover:fill-green-900"} />
                </button>
            </Form>
        </div>
    );
}
