import type {Album} from "~/model/album"
import Preview from "~/components/Preview";
import {formatArtistNames} from "~/model/thing";

type AlbumProps = {
    album: Album;
}

const AlbumPreview = (props: AlbumProps) => {
    const album = props.album;
    const artistNames = formatArtistNames(album.ArtistNames);

    return <Preview artworkLink={album.ArtworkLink} artworkAlt={album.Name}>
        <div className={"text-xl font-bold"}>{album.Name}</div>
        <div>{artistNames}</div>
    </Preview>
}

export default AlbumPreview;
