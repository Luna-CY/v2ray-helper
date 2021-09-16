import {BaseResponse} from "@/api/base"

export const API_V2RAY_SERVER_DEPLOY_STAGE = "/api/v2ray-server-deploy-stage"

export class V2rayServerDeployStageResponse extends BaseResponse {

  public data = new class {

    public stages = new Array<Stage>()

    public cloudreve = new class {

      public email = ""

      public password = ""
    }
  }
}

export class Stage {

  public stage = 0

  public state = 1
}

export class CloudreveInfo {

  public email = ""

  public password = ""
}
