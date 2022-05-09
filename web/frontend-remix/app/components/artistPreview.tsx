import {Artist} from "~/model/artist";
import Preview from "~/components/preview";

type ArtistProps = {
    artist: Artist;
}

const ArtistPreview = (props: ArtistProps) => {
    const artist = props.artist;

    return <Preview artworkLink={artist.ArtworkLink} artworkAlt={artist.Name}>
        <div className={"text-xl font-bold"}>{artist.Name}</div>
    </Preview>
}

export default ArtistPreview;
