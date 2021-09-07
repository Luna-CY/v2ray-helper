/* eslint-disable */
declare module '*.vue' {
  import type {DefineComponent} from 'vue'
  const component: DefineComponent<{}, {}, any>

  export interface ComponentPublicInstance {
    $message: any
  }

  export default component
}
