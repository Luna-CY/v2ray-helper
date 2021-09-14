import {BaseResponse} from "@/api/base"

export const API_META_INFO = "/api/meta-info"

export class MetaInfoResponse extends BaseResponse {

  public data = new class {

    public is_default_key = false

    public is_default_remove_key = false
  }
}
