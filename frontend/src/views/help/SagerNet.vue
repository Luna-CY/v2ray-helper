<template>
  <div class="el-container content-center">
    <div class="el-main el-main-padding hidden-sm-and-down">
      <div class="el-row">
        <div class="el-col-24 helper-box">
          <div class="margin-bottom">第一步：打开APP并点击右上角的"+"，如下图</div>
          <div class="margin-bottom"><img alt="" class="md-and-up-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-ng/v2ray-ng-path-1.png">
          </div>
          <div class="margin-bottom">第二步：选择"扫描二维码"，如下图</div>
          <div class="margin-bottom"><img alt="" class="md-and-up-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-ng/v2ray-ng-path-2.png">
          </div>
          <div class="margin-bottom">第三步：扫描生成的配置二维码后，列表将会出现可用配置项，点击右下角的圆形按钮启动应用，如下图</div>
          <div class="margin-bottom"><img alt="" class="md-and-up-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-ng/v2ray-ng-path-3.png">
          </div>
          <div class="margin-bottom">第四步：点击左上角的"更多"图标，在弹出的菜单中选择设置</div>
          <div class="margin-bottom">第五步：滑动屏幕，找到路由设置下的分应用代理，点击进入</div>
          <div class="margin-bottom">第六步：将上方的状态切换为代理模式，并在下方应用列表中选择需要代理的应用，将其右侧开关打开</div>
          <div class="margin-bottom">第七步：无需保存，退出页面即可</div>
        </div>
      </div>
    </div>
    <div class="el-main el-main-padding hidden-md-and-up">
      <div class="el-row">
        <div class="el-col-24 helper-box">
          <div class="margin-bottom">第一步：打开APP并点击右上角的"+"，如下图</div>
          <div class="margin-bottom"><img alt="" class="sm-and-down-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-ng/v2ray-ng-path-1.png">
          </div>
          <div class="margin-bottom">第二步：选择"扫描二维码"，如下图</div>
          <div class="margin-bottom"><img alt="" class="sm-and-down-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-ng/v2ray-ng-path-2.png">
          </div>
          <div class="margin-bottom">第三步：扫描生成的配置二维码后，列表将会出现可用配置项，点击右下角的圆形按钮启动应用，如下图</div>
          <div class="margin-bottom"><img alt="" class="sm-and-down-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-ng/v2ray-ng-path-3.png">
          </div>
          <div class="margin-bottom">第四步：点击左上角的"更多"图标，在弹出的菜单中选择设置</div>
          <div class="margin-bottom">第五步：滑动屏幕，找到路由设置下的分应用代理，点击进入</div>
          <div class="margin-bottom">第六步：将上方的状态切换为代理模式，并在下方应用列表中选择需要代理的应用，将其右侧开关打开</div>
          <div class="margin-bottom">第七步：无需保存，退出页面即可</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue"

export default defineComponent({
  name: "V2rayNG"
})
</script>

<style lang="scss" scoped>
.el-main.el-main-padding {
  --el-main-padding: 30px;
}

.helper-box {
  background-color: #FFF;
  padding: 30px;
}

.md-and-up-img {
  width: 800px;
}

.sm-and-down-img {
  width: 200px;
}
</style>
