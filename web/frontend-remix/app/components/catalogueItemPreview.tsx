import React, {PropsWithChildren} from "react";


type CatalogueItemPreviewProps = PropsWithChildren<{}> & {
    artworkLink: string;
    artworkAlt: string;
}

const CatalogueItemPreview = (props: CatalogueItemPreviewProps) => {
    return (
        <div className={"flex flex-row gap-1"}>
            <img src={props.artworkLink} alt={props.artworkAlt} className={"h-40 rounded-lg mr-2"}/>
            <div className={"flex flex-col gap-1"}>
                {props.children}
            </div>
        </div>
    );
}

export default CatalogueItemPreview;