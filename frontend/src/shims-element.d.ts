import {IMessage} from "element-plus/lib/el-message/src/types"

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $message: IMessage;
  }
}
