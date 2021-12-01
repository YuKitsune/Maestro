import {PropsWithChildren, useCallback} from "react";

type LinkProps = PropsWithChildren<{}> & {
    link: string
}

const Link = (props: LinkProps) => {
    return <a className={"bg-gray-200 dark:bg-gray-700 rounded-lg p-2 cursor-pointer"} href={props.link} target={"_blank"}>
        {props.children}
    </a>
}

export default Link
