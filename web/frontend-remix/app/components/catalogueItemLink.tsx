import {PropsWithChildren, useCallback} from "react";

type LinkProps = PropsWithChildren<{}> & {
    link: string
}

const CatalogueItemLink = (props: LinkProps) => {
    return <a className={"bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-800 rounded-lg p-2 cursor-pointer"} href={props.link} target={"_blank"}>
        {props.children}
    </a>
}

export default CatalogueItemLink
