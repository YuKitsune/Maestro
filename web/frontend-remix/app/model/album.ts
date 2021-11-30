import { Thing } from "~/model/thing";

export interface Album extends Thing {
    Name        : string;
    ArtistNames : string[];
    ArtworkLink : string;
}
