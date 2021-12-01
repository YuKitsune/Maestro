import type {Album} from "~/model/album"
import {useCallback} from "react";

type AlbumProps = {
    album: Album;
}

const Album = (props: AlbumProps) => {
    const album = props.album;

    const artistNames = album.ArtistNames.join(", ");

    return <div className={"bg-gray-200 dark:bg-gray-700 rounded-lg p-2"}>
        <div className={"flex flex-row gap-1"}>
            <img src={album.ArtworkLink} alt={album.Name} className={"h-40 rounded-lg mr-2"}/>
            <div className={"flex flex-col gap-1"}>
                <div className={"text-xl font-bold"}>{album.Name}</div>
                <div>{artistNames}</div>
            </div>
        </div>
    </div>;
}

export default Album;