export const API_V2RAY_ENDPOINT_NEW = "/api/v2ray-endpoint/new"

export class V2rayEndpointNewForm {

    public cloud = 0

    public endpoint = 0

    public rate = ""

    public remark = ""

    public host = ""

    public port = 443

    public user_id = ""

    public alter_id = 64

    public use_tls = false

    public sni = ""

    public transport_type = 1

    public tcp = new class {

        public type = "none"
    }

    public web_socket = new class {

        public path = ""
    }

    public kcp = new class {

        public type = "none"
    }

    public http2 = new class {

        public host = ""

        public path = ""
    }
}
