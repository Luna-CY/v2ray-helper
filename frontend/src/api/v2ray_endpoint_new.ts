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

  public level = 0

  public transport_type = 1

  public web_socket = new class {

    public path = ""
  }
}
