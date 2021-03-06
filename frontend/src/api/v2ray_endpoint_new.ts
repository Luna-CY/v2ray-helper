import {Header} from "@/api/base"

export const API_V2RAY_ENDPOINT_NEW = "/api/v2ray-endpoint/new"

export class V2rayEndpointNewForm {

  public cloud = 0

  public endpoint = 0

  public rate = ""

  public remark = ""

  public host = ""

  public port = 443

  public user_id = ""

  public alter_id = 0

  public use_tls = false

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

    public path = ""
  }
}
