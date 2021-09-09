import {BaseResponse} from "@/api/base"

export const API_V2RAY_ENDPOINT_DOWNLOAD = "/api/v2ray-endpoint/download"

export class V2rayEndpointDownloadForm {

  public id = 0
}

export class V2rayEndpointDownloadResponse extends BaseResponse {

  public data = new class {

    public content = ""
  }
}
