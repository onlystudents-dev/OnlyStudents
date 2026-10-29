import type {Me} from "./api.tsx";

export type FeatureProps = {
    me: Me,
    fetchMe: FetchMe,
}

export type FetchMe = () => Promise<Me | null>