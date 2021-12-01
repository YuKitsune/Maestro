import type {Track} from "~/model/track"
import {useCallback} from "react";

type TrackProps = {
    track: Track;
}

const Track = (props: TrackProps) => {
    const track = props.track;

    const artistNames = track.ArtistNames.join(", ");

    return <div className={"bg-gray-200 dark:bg-gray-700 rounded-lg p-2"}>
        <div className={"flex flex-row gap-1"}>
            <img src={track.ArtworkLink} alt={track.AlbumName} className={"h-40 rounded-lg mr-2"}/>
            <div className={"flex flex-col gap-1"}>
                <div className={"text-xl font-bold"}>{track.Name}</div>
                <div>{track.AlbumName}</div>
                <div>{artistNames}</div>
            </div>
        </div>
    </div>;
}

export default Track;
