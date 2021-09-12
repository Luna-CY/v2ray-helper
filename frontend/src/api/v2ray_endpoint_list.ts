import {BaseResponse} from "@/api/base"

export const API_V2RAY_ENDPOINT_LIST = "/api/v2ray-endpoint"

export class V2rayEndpointListResponse extends BaseResponse {

    public data = new class {

        public data = new Array<V2rayEndpointListItem>()
    }
}

export class V2rayEndpointListItem {

    public id = 0

    public remark = ""

    public host = ""

    public port = 80

    public user_id = ""

    public alter_id = 4

    public transport_type = 0

    public downloading = false

    public show_delete_button = false

    public show_generate_menu = false
}
