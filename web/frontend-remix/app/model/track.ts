import {Thing} from "~/model/thing";

export interface Track extends Thing {
    Name        :string
    ArtistNames :string[]
    AlbumName   :string
}
