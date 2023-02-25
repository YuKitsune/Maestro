import React from "react";
import Thing from "~/model/thing";
import CardLink from "~/components/cardLink";
import {Service} from "~/model/service";
import CopyLinkButton from "~/components/copyLinkButton";

export type CardProps = {
    renderPreview: () => React.ReactNode;
    items: Thing[];
    services: Service[];
}

const Card = (props: CardProps) => {

    return (
        <div className={"flex flex-col gap-1"}>

            {/* Preview */}
            <div className={"bg-gray-200 dark:bg-gray-700 rounded-lg p-2"}>
                {props.renderPreview()}
            </div>

            <CopyLinkButton />

            {/* Links */}
            <div className={"flex flex-col gap-1"}>
                {props.items.map(t =>
                    <CardLink
                        key={t.Source}
                        link={t.Link}
                        serviceIconUrl={props.services.find(s => s.Key == t.Source)!.LogoURL}
                        serviceName={props.services.find(s => s.Key == t.Source)!.Name} />
                )}
            </div>
        </div>
    );
}

export default Card