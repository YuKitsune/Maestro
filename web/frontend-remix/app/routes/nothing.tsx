import HomeButton from "~/components/HomeButton";
import {MetaFunction} from "remix";

// @ts-ignore
export const meta: MetaFunction = () => {
    return {
        title: "Nothing found...",
        description: "🤦‍"
    };
}

export default function Nothing() {
    return (
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            <div className={"text-2xl text-center"}>
                Sorry, we couldn't find anything...
            </div>
        </div>
    );
}
