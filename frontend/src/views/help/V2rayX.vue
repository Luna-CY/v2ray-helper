<template>
  <div class="el-container content-center">
    <div class="el-main el-main-padding hidden-sm-and-down">
      <div class="el-row">
        <div class="el-col-24 helper-box">
          <div class="margin-bottom">第一步：从顶部菜单找到V2rayX的图标并点击，然后选择"Configure..."，如下图</div>
          <div class="margin-bottom"><img alt="" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-x/v2ray-x-path-1.png"></div>
          <div class="margin-bottom">第二步：在弹出的窗口中点击"Import"按钮，并从下拉菜单中选择"Import from other links"选项，如下图</div>
          <div class="margin-bottom"><img alt="" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-x/v2ray-x-path-3.png"></div>
          <div class="margin-bottom">第三步：将生成的配置链接粘贴入输入框中（如果连接末尾有等号，注意不要遗漏链接最后的等号），如下图</div>
          <div class="margin-bottom"><img alt="" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-x/v2ray-x-path-4.png"></div>
          <div class="margin-bottom">第四步：从顶部菜单找到V2rayX的图标并点击，然后选择"Load core"，如下图</div>
          <div class="margin-bottom"><img alt="" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-x/v2ray-x-path-5.png"></div>
        </div>
      </div>
    </div>
    <div class="el-main el-main-padding hidden-md-and-up">
      <div class="el-row">
        <div class="el-col-24 helper-box">
          <div class="margin-bottom">第一步：从顶部菜单找到V2rayX的图标并点击，然后选择"Configure..."，如下图</div>
          <div class="margin-bottom"><img alt="" class="sm-and-down-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-x/v2ray-x-path-1.png">
          </div>
          <div class="margin-bottom">第二步：在弹出的窗口中点击"Import"按钮，并从下拉菜单中选择"Import from other links"选项，如下图</div>
          <div class="margin-bottom"><img alt="" class="sm-and-down-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-x/v2ray-x-path-3.png">
          </div>
          <div class="margin-bottom">第三步：将生成的配置链接粘贴入输入框中（如果连接末尾有等号，注意不要遗漏链接最后的等号），如下图</div>
          <div class="margin-bottom"><img alt="" class="sm-and-down-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-x/v2ray-x-path-4.png">
          </div>
          <div class="margin-bottom">第四步：从顶部菜单找到V2rayX的图标并点击，然后选择"Load core"，如下图</div>
          <div class="margin-bottom"><img alt="" class="sm-and-down-img" src="https://github.com/Luna-CY/v2ray-helper/raw/master/web/img/help/v2ray-x/v2ray-x-path-5.png">
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue"

export default defineComponent({
  name: "V2rayX"
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

.sm-and-down-img {
  width: 120px;
}
</style>
