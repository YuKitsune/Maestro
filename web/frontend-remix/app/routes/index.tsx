import type {MetaFunction, ActionFunction} from "remix";
import {redirect, Form} from "remix";

// https://remix.run/api/conventions#meta
export let meta: MetaFunction = () => {
  return {
    title: "Remix Starter",
    description: "Welcome to remix!"
  };
};

export let action: ActionFunction = async ({ request }) => {
  let formData = await request.formData();

  let link = formData.get("link") as string;

  link = encodeURIComponent(link);

  const apiUrl = `${process.env.API_URL}/link?link=${link}`
  const res = await fetch(apiUrl)
  const json = await res.json()

  // Todo: Get group ID
  const groupId = json[0].GroupId;

  return redirect(`/${groupId}`);
};

// https://remix.run/guides/routing#index-routes
export default function Index() {
  return (
    <div className={"flex flex-col gap-1 align-center items-center"}>
      <h1 className={"text-4xl text-align-center mb-2"}>Welcome to Maestro!</h1>
      <h3 className={"text-xl"}>Paste a link to an Artist, Album or Track here:</h3>
      <Form method={"post"} className={"w-full flex flex-row gap-2"}>
        <input
            className={"w-full rounded-lg border-2 border-blue-100 dark:border-blue-900 dark:bg-black px-1 focus-within:border-blue-500 outline-none"}
            type="text"
            placeholder={"https://open.spotify.com/track/4cOdK2wGLETKBW3PvgPWqT"}
            name={"link"}/>
        <button type="submit" className={"rounded-lg px-1 bg-green-200 dark:bg-green-800"}>Go!</button>
      </Form>
    </div>
  );
}