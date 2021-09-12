export const API_V2RAY_SERVER_DEPLOY = "/api/v2ray-server-deploy"

export class V2rayServerDeployForm {

    // 服务器信息

    public server_type = 1 // 1本地服务器；2远程服务器

    public server_ip = ""

    public server_port = 22

    public server_user = ""

    public server_password = ""

    public install_type = 1 // 1默认安装；2强制安装；3升级安装；4重新配置；5删除V2ray

    // V2ray部署配置

    public config_type = 1 // 1预设配置(WebSocket/HTTP)；2预设配置(WebSocket/HTTPS)；3自定义配置

    public clients = new Array<Client>()

    public v2ray_port = 3000

    public transport_type = 1

    public tcp = new class {

        public type = "none"

        public request = new class {

            public version = "1.1"

            public method = "GET"

            public path = ""

            public headers = new Array<Header>()
        }

        public response = new class {

            public version = "1.1"

            public status = 200

            public reason = "OK"

            public headers = new Array<Header>()
        }
    }

    public web_socket = new class {

        public path = ""

        public headers = new Array<Header>()
    }

    public kcp = new class {

        public mtu = 1350

        public tti = 50

        public uplink_capacity = 5

        public downlink_capacity = 20

        public congestion = false

        public read_buffer_size = 2

        public write_buffer_size = 2

        public type = "none"
    }

    public http2 = new class {

        public host = ""

        public path = ""
    }

    // 其他配置

    public use_tls = false

    public tls_host = ""

    public cert_type = 1 // 证书类型：1申请新证书；2上传证书
}

export class Client {

    public user_id = ""

    public alter_id = 4
}

export class Header {

    public key = ""

    public value = ""
}
