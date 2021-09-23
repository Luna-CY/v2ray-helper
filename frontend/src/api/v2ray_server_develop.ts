import {BaseResponse, Header} from "@/api/base"

export const API_V2RAY_SERVER_DEPLOY = "/api/v2ray-server-deploy"

export class V2rayServerDeployForm {

  public install_type = 1 // 1默认安装；2强制安装；3升级安装；4重新配置

  public management_key = "" // 管理员口令

  // HTTPS配置

  public use_tls = false

  public tls_host = ""

  public cert_type = 1 // 证书类型：1申请新证书；2上传证书

  // 站点伪装配置

  public enable_web_service = false

  public web_service_type = "cloudreve"

  public cloudreve_config = new class {

    public enable_aria2 = false

    public reset_admin_password = true
  }

  // V2ray部署配置

  public config_type = 1 // 1预设配置(WebSocket/HTTP)；2预设配置(WebSocket/HTTPS)；3自定义配置

  public v2ray_config = new class {

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

        public status = "200"

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

      public path = "/"
    }
  }
}

export class V2rayServerDeployResponse extends BaseResponse {

  public data = new V2rayServerDeployData()
}

export class V2rayServerDeployData {

  public cloudreve_admin = ""

  public cloudreve_password = ""
}

export class Client {

  public user_id = ""

  public alter_id = 4
}
