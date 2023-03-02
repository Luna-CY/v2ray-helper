import {Response} from "@/apis/types";
import axios from "@/apis/axios";

export declare type GuardBody = { key: string }
export declare type GuardResponse = { token: string, expired: number }

const api = "/auth"

export const guard = async (body: GuardBody): Promise<Response<GuardResponse>> => {
    const response = await axios.post(api, body)

    if (200 !== response.status || 0 !== response.data.code) {
        return {code: response.data.code, message: response.data.message}
    }

    return {code: response.data.code, message: response.data.message, data: response.data.data}
}
