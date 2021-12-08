import React from "react";
import QuestionMarkIcon from "~/components/icons/questionMark";
import {Thing} from "~/model/thing";
import {Service} from "~/model/service";

type LinkProps = {
    thing: Thing
    service: Service | undefined
}

const CatalogueItemLink = ({thing, service}: LinkProps) => {

    const icon = service !== undefined ?
        <img src={service.ArtworkLink} alt={thing.Source} className={"row-span-2 flex-shrink h-12"}/> :
        <QuestionMarkIcon className={"h-12 w-12"} />;

    return <a className={"bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-800 rounded-lg p-2 cursor-pointer"} href={thing.Link} target={"_blank"}>
        <div className={"flex flex-row content-center gap-1"}>

            <div className={"w-12"}>
                {icon}
            </div>

            <div className={"grid"}>
                <div className={"font-bold"}>
                    {service?.Name ?? "Unknown"}
                </div>
                <div className={"text-underline text-blue-400 truncate"}>
                    {thing.Link}
                </div>
            </div>
        </div>
    </a>
}

export default CatalogueItemLink
