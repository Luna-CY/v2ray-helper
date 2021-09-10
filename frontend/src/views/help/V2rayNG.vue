<template>
  <div class="el-container content-center">
    <div class="el-main el-main-padding hidden-sm-and-down">
      <div class="el-row">
        <div class="el-col-24 helper-box">
          <div class="margin-bottom">第一步：打开APP并点击右上角的"+"，如下图</div>
          <div class="margin-bottom"><img class="md-and-up-img" src="/img/help/v2ray-ng/v2ray-ng-path-1.png" alt=""></div>
          <div class="margin-bottom">第二步：选择"扫描二维码"，如下图</div>
          <div class="margin-bottom"><img class="md-and-up-img" src="/img/help/v2ray-ng/v2ray-ng-path-2.png" alt=""></div>
          <div class="margin-bottom">第三步：扫描生成的配置二维码后，列表将会出现可用配置项，点击右下角的圆形按钮启动应用，如下图</div>
          <div class="margin-bottom"><img class="md-and-up-img" src="/img/help/v2ray-ng/v2ray-ng-path-3.png" alt=""></div>
        </div>
      </div>
    </div>
    <div class="el-main el-main-padding hidden-md-and-up">
      <div class="el-row">
        <div class="el-col-24 helper-box">
          <div class="margin-bottom">第一步：打开APP并点击右上角的"+"，如下图</div>
          <div class="margin-bottom"><img class="sm-and-down-img" src="/img/help/v2ray-ng/v2ray-ng-path-1.png" alt=""></div>
          <div class="margin-bottom">第二步：选择"扫描二维码"，如下图</div>
          <div class="margin-bottom"><img class="sm-and-down-img" src="/img/help/v2ray-ng/v2ray-ng-path-2.png" alt=""></div>
          <div class="margin-bottom">第三步：扫描生成的配置二维码后，列表将会出现可用配置项，点击右下角的圆形按钮启动应用，如下图</div>
          <div class="margin-bottom"><img class="sm-and-down-img" src="/img/help/v2ray-ng/v2ray-ng-path-3.png" alt=""></div>
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

<style scoped lang="scss">
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
