import React from "react";

type LinkProps = {
    link: string
    serviceIconUrl: string
    serviceName: string
}

const CardLink = (props: LinkProps) => {
    return <a className={"bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-800 rounded-lg p-2 cursor-pointer"} href={props.link} target={"_blank"}>
        <div className={"flex flex-row content-center gap-2"}>

            <div className={"w-12 h-12 flex-none"}>
                <img src={props.serviceIconUrl} alt={props.serviceName} className={"row-span-2 flex-shrink h-12 w-12"}/>
            </div>

            <div className={"grid"}>
                <div className={"font-bold"}>
                    {props.serviceName ?? "Unknown"}
                </div>
                <div className={"text-underline text-blue-400 truncate"}>
                    {props.link}
                </div>
            </div>
        </div>
    </a>
}

export default CardLink