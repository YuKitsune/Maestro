import HomeButton from "~/components/homeButton";

export default function Group() {
    return (
        <div className={"flex flex-col gap-2"}>
            <HomeButton />
            <div className={"text-2xl text-center"}>
                Sorry, we couldn't find anything...
            </div>
        </div>
    );
}
