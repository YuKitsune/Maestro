import {
    Links,
    LiveReload,
    Meta,
    Outlet,
    Scripts,
    ScrollRestoration, useCatch,
} from "remix";
import type { LinksFunction } from "remix";

import tailwindStylesUrl from "./tailwind.css"
import spinnerStylesUrl from "./components/spinner.css"
import HomeButton from "~/components/homeButton";

// https://remix.run/api/app#links
export let links: LinksFunction = () => {
  return [
    { rel: "stylesheet", href: tailwindStylesUrl },
    { rel: "stylesheet", href: spinnerStylesUrl },
  ];
};

// https://remix.run/api/conventions#default-export
// https://remix.run/api/conventions#route-filenames
export default function App() {
  return (
    <Document>
      <Layout>
        <Outlet />
      </Layout>
    </Document>
  );
}

// https://remix.run/docs/en/v1/api/conventions#errorboundary
export function ErrorBoundary({ error }: { error: Error }) {
  console.error(error);
  return (
    <Document title="Error!">
      <Layout>
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            <h1 className={"text-2xl text-center"}>Something went wrong 😔</h1>
            {error.message && <div className={"bg-gray-200 dark:bg-gray-700 rounded-lg p-2 mb-1 font-mono"}>{error.message}</div>}
            <div className={"text-center"}>
                <p>Try again and see if it works.</p>
                <p>Feel free to open a <a className={"underline text-blue-400"} href={"https://github.com/YuKitsune/Maestro/issues/new"}>GitHub issue</a>.</p>
            </div>
        </div>
      </Layout>
    </Document>
  );
}

// https://remix.run/docs/en/v1/api/conventions#catchboundary
export function CatchBoundary() {
  let caught = useCatch();

  let message;
  switch (caught.status) {
    case 401:
      message = <span>Huh. You shouldn't be here... 👮</span>
      break;

    case 404:
        message = <span>Where are you going? 🤔</span>
      break;

    default:
      throw new Error(caught.data || caught.statusText);
  }

  return (
    <Document title={`${caught.status} ${caught.statusText}`}>
      <Layout>
          <div className={"flex flex-col gap-2"}>
              <HomeButton />
              <h1 className={"text-2xl text-center"}>{message}</h1>
          </div>
      </Layout>
    </Document>
  );
}

function Document({
  children,
  title
}: {
  children: React.ReactNode;
  title?: string;
}) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width,initial-scale=1" />

        {/* Emoji favicon because I can't do graphic design */}
        <link rel="icon" href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>🎵</text></svg>" />

        {title ? <title>{title}</title> : null}
        <Meta />
        <Links />
      </head>
      <body>
        {children}
        <ScrollRestoration />
        <Scripts />
        {process.env.NODE_ENV === "development" && <LiveReload />}
      </body>
    </html>
  );
}

function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex flex-col flex-grow items-center justify-between h-screen overflow-y-scroll p-8 gap-8 bg-gray-200 dark:bg-gray-800 dark:text-white">
      <div className={"flex-grow rounded-lg bg-gray-100 dark:bg-gray-900 shadow-2xl p-4 w-full max-w-md"}>
        {children}
      </div>
      <footer className={"flex flex-col items-center"}>
        <a href="https://www.buymeacoffee.com/yukitsune256" target="_blank">
          <img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" className={"w-40 pb-2"} />
        </a>
        <div>
          <p><a className={"underline text-blue-400"} href={"https://github.com/YuKitsune/Maestro"}>GitHub</a></p>
        </div>
      </footer>
    </div>
  );
}
