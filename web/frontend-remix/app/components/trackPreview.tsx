import type {Track} from "~/model/track"
import CatalogueItemPreview from "~/components/catalogueItemPreview";

type TrackProps = {
    track: Track;
}

const TrackPreview = (props: TrackProps) => {
    const track = props.track;
    const artistNames = track.ArtistNames.join(", ");

    return <CatalogueItemPreview artworkLink={track.ArtworkLink} artworkAlt={track.Name}>
            <div className={"text-xl font-bold"}>{track.Name}</div>
            <div>{artistNames}</div>
    </CatalogueItemPreview>
}

export default TrackPreview;
