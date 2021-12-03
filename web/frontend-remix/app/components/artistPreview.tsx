import {Artist} from "~/model/artist";
import CatalogueItemPreview from "~/components/catalogueItemPreview";

type ArtistProps = {
    artist: Artist;
}

const ArtistPreview = (props: ArtistProps) => {
    const artist = props.artist;

    return <CatalogueItemPreview artworkLink={artist.ArtworkLink} artworkAlt={artist.Name}>
            <div className={"text-xl font-bold"}>{artist.Name}</div>
    </CatalogueItemPreview>
}

export default ArtistPreview;
