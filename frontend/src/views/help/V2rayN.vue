<template>
  <div class="el-container content-center">
    <div class="el-main el-main-padding">
      <div class="el-row">
        <div class="el-col-24 helper-box">
          <div>暂时没有教程</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue"

export default defineComponent({
  name: "V2rayN"
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
</style>
