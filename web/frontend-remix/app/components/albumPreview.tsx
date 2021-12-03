import type {Album} from "~/model/album"
import CatalogueItemPreview from "~/components/catalogueItemPreview";

type AlbumProps = {
    album: Album;
}

const AlbumPreview = (props: AlbumProps) => {
    const album = props.album;
    const artistNames = album.ArtistNames.join(", ");

    return <CatalogueItemPreview artworkLink={album.ArtworkLink} artworkAlt={album.Name}>
        <div className={"text-xl font-bold"}>{album.Name}</div>
        <div>{artistNames}</div>
    </CatalogueItemPreview>
}

export default AlbumPreview;
