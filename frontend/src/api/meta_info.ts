import {BaseResponse} from "@/api/base"

export const API_META_INFO = "/api/meta-info"

export class MetaInfoResponse extends BaseResponse {

  public data = new class {

    public is_default_key = false

    public is_default_remove_key = false

    public listen = 8888

    public enable_https = false

    public https_host = ""

    public email = ""

    public notice_list = new Array<NoticeListItem>()
  }
}

export class NoticeListItem {

  public id = ""

  public type = 1

  public title = ""

  public message = ""

  public time = 0
}
