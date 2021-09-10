export const API_V2RAY_SERVER_DEPLOY = "/api/v2ray-server-deploy"

export class V2rayServerDeployForm {

  // 服务器信息

  public server_is_local = false

  public server_ip = ""

  public server_port = 22

  public server_user = ""

  public server_password = ""

  // V2ray部署配置

  public alter_id = 64

  public v2ray_port = 3000

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

  // 其他配置

  public use_tls = false

  public tls_host = ""

  public use_nginx_proxy = false
}
