import {Error as ErrorModel} from "~/model/error"
import {Service} from "~/model/service";
import {Artist} from "~/model/artist";
import {Album} from "~/model/album";
import {Track} from "~/model/track";

export interface Response<T> {
    Type: "artist" | "album" | "track";
    Items: T[]
}

class MaestroApiClient {
    readonly baseUrl: string;
    readonly publicBaseUrl: string;

    constructor(baseUrl: string, publicBaseUrl: string) {
        this.baseUrl = baseUrl;
        this.publicBaseUrl = publicBaseUrl;
    }

    async searchFromLink(link: string): Promise<Response<any>> {
        const urlSafeLink = encodeURIComponent(link);
        const apiUrl = `${this.baseUrl}/link?link=${urlSafeLink}`

        const res = await fetch(apiUrl)
        const json = await res.json()

        if (res.status != 200) {
            const resErr = json as ErrorModel
            throw new Error(resErr.Error)
        }

        return json as Response<any>;
    }

    async getArtists(id: string): Promise<Artist[]> {
        const res = await this.tryGetArtists(id);
        const err = res as ErrorModel

        if (err.Error) {
            throw new Error(err.Error)
        }

        return res as Artist[];
    }

    async tryGetArtists(id: string): Promise<Artist[] | ErrorModel> {
        const apiUrl = `${this.baseUrl}/artist/${id}`

        const res = await fetch(apiUrl)
        const json = await res.json()

        if (res.status != 200) {
            const resErr = json as ErrorModel
            return resErr
        }

        const artistRes = json as Response<Artist>;
        return artistRes.Items;
    }

    async getAlbums(id: string): Promise<Album[]> {
        const res = await this.tryGetAlbums(id);
        const err = res as ErrorModel

        if (err.Error) {
            throw new Error(err.Error)
        }

        return res as Album[];
    }

    async tryGetAlbums(id: string): Promise<Album[] | ErrorModel> {
        const apiUrl = `${this.baseUrl}/album/${id}`

        const res = await fetch(apiUrl)
        const json = await res.json()

        if (res.status != 200) {
            const resErr = json as ErrorModel
            return resErr;
        }

        const albumRes = json as Response<Album>;
        return albumRes.Items;
    }

    async getTracks(isrc: string): Promise<Track[]> {
        const res = await this.tryGetTracks(isrc);
        const err = res as ErrorModel

        if (err.Error) {
            throw new Error(err.Error)
        }

        return res as Track[];
    }

    async tryGetTracks(isrc: string): Promise<Track[] | ErrorModel> {
        const apiUrl = `${this.baseUrl}/track/${isrc}`

        const res = await fetch(apiUrl)
        const json = await res.json()

        if (res.status != 200) {
            const resErr = json as ErrorModel
            return resErr;
        }

        const trackRes = json as Response<Track>;
        return trackRes.Items;
    }

    async getServices(): Promise<Service[]> {
        const apiUrl = `${this.baseUrl}/services`

        const res = await fetch(apiUrl)
        const json = await res.json()

        if (res.status != 200) {
            const resErr = json as ErrorModel
            throw new Error(resErr.Error)
        }

        let services = json as Service[];

        // Get artwork URLs
        for (const svc of services) {
            // Todo: Need a public URL here, the docker compose one isn't accessible publicly
            svc.LogoURL = `${this.publicBaseUrl}/services/${svc.Key}/logo`;
        }

        return services;
    }
}

export default MaestroApiClient;