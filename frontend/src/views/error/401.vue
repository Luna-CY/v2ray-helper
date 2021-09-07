<template>
  <div class="error-401">
    <div class="content-center margin-bottom-2x"><i class="el-icon-lock icon-401"></i></div>
    <div class="content-center">
      <span>您没有权限访问该页面</span>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'

export default defineComponent({
  name: "401"
})
</script>

<style lang="scss" scoped>
.error-401 {
  margin: 240px auto 0 auto;
}

.icon-401 {
  font-size: 8em;
}
</style>
