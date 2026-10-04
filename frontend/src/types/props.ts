import type {Me} from "./api.ts";

export type FeatureProps = {
    me: Me,
    unauthorized: () => Promise<Me | null>,
}