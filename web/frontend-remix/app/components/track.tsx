import type {Track} from "~/model/track"
import {useCallback} from "react";

type TrackProps = {
    track: Track;
}

const Track = (props: TrackProps) => {
    const track = props.track;

    const artistNames = track.ArtistNames.join(", ");

    const openLink = useCallback(() => {
        window.open(track.Link, "_blank");
    }, [props])

    return <div className={"bg-gray-200 rounded-lg p-2 cursor-pointer"} onClick={openLink}>
        <div>{track.Name}</div>
        <div>{artistNames}</div>
        <div>{track.AlbumName}</div>
        <div>{track.Source}</div>
    </div>;
}

export default Track;
