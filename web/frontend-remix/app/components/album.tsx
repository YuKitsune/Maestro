import type {Album} from "~/model/album"
import {useCallback} from "react";

type AlbumProps = {
    album: Album;
}

const Album = (props: AlbumProps) => {
    const album = props.album;

    const artistNames = album.ArtistNames.join(", ");

    const openLink = useCallback(() => {
        window.open(album.Link, "_blank");
    }, [props])

    return <div className={"bg-gray-200 rounded-lg p-2 cursor-pointer"} onClick={openLink}>
        <div>{album.Name}</div>
        <div>{artistNames}</div>
        <div>{album.Source}</div>
    </div>;
}

export default Album;