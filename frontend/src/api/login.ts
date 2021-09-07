import {BaseResponse} from "@/api/base"

export const API_LOGIN = "/api/auth"

export class LoginForm {

  public key: string = ''
}

export class LoginResponse extends BaseResponse {

  public data = new class {

    public token: string = ''

    public expired: number = 0
  }
}
