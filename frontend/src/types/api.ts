import type {Theme} from "../util/theme.ts";

export type Role = "student" | "teacher" | "guardian" | "admin";

export type Me = {
    role: Role,
    account_id: number,
    first_name: string,
    last_name: string,
    pfp_url: string,
    email_address: string,
    email_verified: boolean,
    nickname: string,
    lang: string,
    theme: Theme,
}
