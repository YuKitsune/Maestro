import {Artist} from "~/model/artist";
import {useCallback} from "react";

type ArtistProps = {
    artist: Artist;
}

const Artist = (props: ArtistProps) => {

    const artist = props.artist;

    return <div className={"bg-gray-200 dark:bg-gray-700 rounded-lg p-2"}>
        <div className={"flex flex-row gap-1"}>
            <img src={artist.ArtworkLink} alt={artist.Name} className={"h-40 rounded-lg mr-2"}/>
            <div className={"flex flex-col gap-1"}>
                <div className={"text-xl font-bold"}>{artist.Name}</div>
            </div>
        </div>
    </div>;
}

export default Artist;
