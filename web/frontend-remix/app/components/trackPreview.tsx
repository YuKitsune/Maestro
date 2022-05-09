import type {Track} from "~/model/track"
import Preview from "~/components/preview";
import {formatArtistNames} from "~/model/thing";

type TrackProps = {
    track: Track;
}

const TrackPreview = (props: TrackProps) => {
    const track = props.track;
    const artistNames = formatArtistNames(track.ArtistNames);

    return <Preview artworkLink={track.ArtworkLink} artworkAlt={track.Name}>
        <div className={"text-xl font-bold"}>{track.Name}</div>
        <div>{artistNames}</div>
    </Preview>
}

export default TrackPreview;
