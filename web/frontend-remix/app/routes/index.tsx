import type {ActionFunction, MetaFunction} from "remix";
import {Form, redirect, useTransition} from "remix";
import {Thing} from "~/model/thing";
import Spinner from "~/components/spinner";

// https://remix.run/api/conventions#meta
export let meta: MetaFunction = () => {
  return {
    title: "Maestro",
    description: "Share links to many popular music streaming services!"
  };
};

export const action: ActionFunction = async ({request}) => {

    const formData = await request.formData();
    const link = formData.get("link") as string;
    const urlSafeLink = encodeURIComponent(link);
    const apiUrl = `${process.env.API_URL}/link?link=${urlSafeLink}`
    const res = await fetch(apiUrl)
    const json = await res.json()

    const things = json as Thing[];

    if (things && things.length > 0) {
        const groupId = things[0].GroupId;
        return redirect(`/${groupId}`)
    }

    // Todo: Navigate to a "No results page"
}

// https://remix.run/guides/routing#index-routes
export default function Index() {
    const transition = useTransition();

    let isLoading = transition.state !== "idle";

    return (
        <div className={"flex flex-col gap-2 align-center items-center"}>
            <h1 className={"text-4xl text-align-center"}>Welcome to Maestro!</h1>
            {!isLoading && (
                <>
                    <div className={"text-align-center"}>Paste a link to an Artist, Album or Track here:</div>
                    <Form method={"post"} className={"w-full flex flex-row gap-2"}>
                        <input
                            className={"w-full rounded-lg border-2 border-blue-100 dark:border-blue-900 dark:bg-black px-1 focus-within:border-blue-500 outline-none"}
                            type="text"
                            placeholder={"https://open.spotify.com/track/4cOdK2wGLETKBW3PvgPWqT"}
                            name={"link"}
                        />
                        <button type="submit" className={"rounded-lg px-1 bg-green-300 hover:bg-green-400 dark:bg-green-800 dark:hover:bg-green-900"}>Go!</button>
                    </Form>
                </>
            )}

            {isLoading && (
                <>
                    {/* Todo: Randomize the loading messages */}
                    <div className={"text-align-center"}>Give us a sec...</div>
                    <Spinner />
                </>
            )}
        </div>
    );
}
