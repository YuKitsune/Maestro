import {Artist} from "~/model/artist";
import {useCallback} from "react";

type ArtistProps = {
    artist: Artist;
}

const Artist = (props: ArtistProps) => {

    const artist = props.artist;

    const openLink = useCallback(() => {
        window.open(artist.Link, "_blank");
    }, [props])

    return <div className={"bg-gray-200 rounded-lg p-2 cursor-pointer"} onClick={openLink}>
        <div>{artist.Name}</div>
        <div>{artist.Source}</div>
    </div>;
}

export default Artist;
